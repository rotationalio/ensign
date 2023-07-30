#!/usr/bin/env python3
# Uses gcloud commands to cleanup old GCR images and reduce storage costs.

import re
import json
import argparse
import subprocess
from datetime import datetime, timedelta, timezone


semver = re.compile(r'^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$', re.I)
digit2 = re.compile(r'v?(0|[1-9]\d*)\.(0|[1-9]\d*)', re.I)
latest = re.compile(r'\s*latest\s*', re.I)


def get_images(repo=""):
    """
    Returns the images managed by GCR for the specified repository.
    """
    result = subprocess.run([
        "gcloud", "container", "images", "list",
        "--repository", repo,
        "--format", "json",
    ], capture_output=True)

    if result.returncode > 0:
        raise Exception(result.stderr.decode('utf-8').strip())

    data = json.loads(result.stdout)
    return [row['name'] for row in data]


def get_digests(image):
    """
    Returns all the digests being managed for the specified image.
    """
    result = subprocess.run([
        "gcloud", "container", "images", "list-tags", image,
        "--limit", "unlimited",
        "--sort-by", "timestamp",
        "--format", "json"
    ], capture_output=True)

    if result.returncode > 0:
        raise Exception(result.stderr.decode('utf-8').strip())
    return json.loads(result.stdout)


def filter_digests(digests, n_keep, grace):
    """
    Filters digests and returns a list of digests to delete and a list of digests to
    keep. The filtering process is as follows:
    1. If the digest has no tag associated with it, delete it
    2. If the digest has a tag that matches a semantic version, then keep it
    3. If the digest is within the grace period, then keep it
    4. All other digests are moved to delete
    5. After processing move the number of keep digests from delete to keep so long as
       they have a tag (e.g. all images without tags are deleted).
    """
    now = datetime.now(timezone.utc)
    grace = timedelta(hours=grace)

    n_grace = 0
    keep, delete = [], []
    for digest in digests:
        # Delete any digests that have no tags
        if len(digest['tags']) == 0:
            delete.append(digest)
            continue

        # Check if the digest has a semver or latest tag
        for tag in digest['tags']:
            if semver.match(tag) or digit2.match(tag) or latest.match(tag):
                keep.append(digest)
                break
        else:
            # The for loop did not break so the digest was not appended to keep
            # Check if it's in the grace period otherwise delete it
            ts = datetime.strptime(digest["timestamp"]["datetime"], "%Y-%m-%d %H:%M:%S%z")
            if now - ts < grace:
                keep.append(digest)
                n_grace += 1
            else:
                delete.append(digest)

    # Sort deleted by timestamp so only the most recent are kept
    delete.sort(reverse=True, key=lambda d: d['timestamp']['datetime'])

    if n_grace < n_keep:
        # Move items from delete as long as it is tagged until we have n_keep items
        candidates = []
        for i, digest in enumerate(delete):
            if len(digest['tags']) > 0:
                candidates.append((i, digest))

            if len(candidates) >= (n_keep - n_grace):
                break

        seti = set([])
        for i, digest in candidates:
            keep.append(digest)
            seti.add(i)
        delete = [digest for i, digest in enumerate(delete) if i not in seti]

    # Sort keep by timestamp so we know what's being kept
    keep.sort(reverse=True, key=lambda d: d['timestamp']['datetime'])
    return keep, delete


def delete_image(image, digest):
    result = subprocess.run([
        "gcloud", "container", "images", "delete",
        f"{image}@{digest}",
        "--force-delete-tags",
        "--quiet",
    ], capture_output=True)

    if result.returncode > 0:
        raise Exception(result.stderr.decode('utf-8').strip())
    cprint(f"deleted {image}@{digest}", rgb.RED)


def prompt(msg):
    while True:
        answer = input(f"{msg} [Y|n]: ").lower()
        if answer[0] == 'y':
            return True
        if answer[0] == 'n':
            return False


def cprint(text, rgb):
    print(cfmt(text, rgb))


def cfmt(text, rgb):
    r, g, b = rgb
    return f"\033[38;2;{r};{g};{b}m{text}\033[0m"


class rgb():
    BLACK = (0, 0, 0)
    RED = (255, 0, 0)
    GREEN = (0, 255, 0)
    BLUE = (0, 0, 255)
    YELLOW = (255, 255, 0)


def main(args):
    cprint("Getting all images for repository", rgb.YELLOW)
    images = get_images(args.repo)
    print("\n".join(images))

    for image in images:
        cprint(f"Checking {image} for cleanup requirements", rgb.YELLOW)
        digests = get_digests(image)
        keep, delete = filter_digests(digests, args.keep, args.grace)

        if args.dryrun:
            if len(delete) == 0:
                cprint("no images are candidates for deletion", rgb.GREEN)
            else:
                nkeep = cfmt(f"{len(keep)}", rgb.GREEN)
                ndel = cfmt(f"{len(delete)}", rgb.RED)
                print(f"{ndel} images to delete {nkeep} images to keep")

                if args.yes or prompt(cfmt(f"display {len(delete)} images being deleted?", rgb.YELLOW)):
                   for d in delete:
                        cprint(f"{d['digest']} {d['timestamp']['datetime']} {d['tags']}", rgb.RED)

            if args.yes or prompt(cfmt(f"display {len(keep)} images being kept?", rgb.YELLOW)):
                for d in keep:
                    print(f"{d['digest']} {d['timestamp']['datetime']} {d['tags']}")

            continue

        else:
            if len(delete) == 0:
                nkeep = cfmt(f"{len(keep)}", rgb.GREEN)
                print(f"no images are candidates for deletion, keeping {nkeep} images")
            else:
                nkeep = cfmt(f"{len(keep)}", rgb.GREEN)
                ndel = cfmt(f"{len(delete)}", rgb.RED)
                print(f"{ndel} images to delete {nkeep} images to keep")

                n_deleted = 0
                if args.yes or prompt(cfmt(f"delete {len(delete)} images?", rgb.YELLOW)):
                    # Reverse the order of the images to delete to reduce dependencies
                    delete.sort(key=lambda d:d['timestamp']['datetime'])

                    for digest in delete:
                        try:
                            delete_image(image, digest["digest"])
                            n_deleted += 1
                        except Exception as e:
                            print(e)
                            continue

                if n_deleted != len(delete):
                    cprint(f"deleted {n_deleted} out of {len(delete)} images", rgb.RED)


if __name__ == "__main__":
    args = {
        ("-r", "--repo"): {
            "type": str, "required": False, "default":"",
            "help": "specify the repository to list images for, leave blank to use the current project",
        },
        ("-k", "--keep"): {
            "type": int, "default": 3, "metavar": "N",
            "help": "specify the minimum number of images to keep",
        },
        ("-g", "--grace"): {
            "type": int, "default": 168, "metavar": "H",
            "help": "minimum duration (in hours) to ignore references, e.g. images more recent than the duration will be kept",
        },
        ("-d", "--dryrun"): {
            "action": "store_true", "default": False,
            "help": "show what will be deleted without actually deleting anything",
        },
        ("-y", "--yes"): {
            "action": "store_true", "default": False,
            "help": "automatically respond yes to all prompts",
        },
        ("-T", "--traceback"): {
            "action": "store_true", "default": False,
            "help": "print stack trace on error",
        }
    }

    parser = argparse.ArgumentParser(
        description="Uses gcloud commands to cleanup old GCR images and reduce storage costs.",
        epilog="Note you must be logged into gcloud with access to the project before running this script",
    )

    for pargs, kwargs in args.items():
        if isinstance(pargs, str):
            pargs = (pargs )
        parser.add_argument(*pargs, **kwargs)

    args = parser.parse_args()
    try:
        main(args)
    except Exception as e:
        if args.traceback:
            raise e
        else:
            parser.error(e)
    except KeyboardInterrupt:
        print()
        parser.error("operation canceled by user")