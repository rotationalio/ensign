import slugify from "../slugifyDomain";
import { describe, expect, it } from "vitest";

describe("#slugify", () => {
    it("returns ensign.rotational.io when org name is empty", () => {
        expect(slugify("")).toBe("ensign.rotational.io/")
    })

    it("returns ensign.rotational.io/rotational-labs-inc when org name is Roational Labs, Inc", () => {
        expect(slugify("Rotational Labs, Inc.")).toBe("ensign.rotational.io/rotational-labs-inc")
    })

    it("returns ensign.rotational.io/hermes-international-sa when org name is Hermès International S.A.", () => {
        expect(slugify("Hermès International S.A.")).toBe("ensign.rotational.io/hermes-international-sa")
    })

    it("returns ensign.rotational.io/baskin-robins when org name is Baskin Robins", () => {
        expect(slugify("Baskin-Robins")).toBe("ensign.rotational.io/baskin-robins")
    })
})