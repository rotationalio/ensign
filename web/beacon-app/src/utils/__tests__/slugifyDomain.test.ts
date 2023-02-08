/* eslint-disable prettier/prettier */
import { describe, expect, it } from 'vitest';

import { slugify } from '../slugifyDomain';

describe('#slugify', () => {
  it('returns ensign.rotational.io when org name is empty', () => {
    expect(slugify('')).toBe('ensign.rotational.io/');
  });

  it('returns ensign.rotational.io/rotational-labs-inc/domain-space when org name is Roational Labs, Inc.', () => {
    expect(slugify('domain space', 'Rotational Labs, Inc.')).toBe(
      'ensign.rotational.io/rotational-labs-inc/domain-space'
    );
  });

  // it("returns ensign.rotational.io/hermes-international-sa when org name is Hermès International S.A.", () => {
  //     expect(slugify("Hermès International S.A.", "my org")).toBe("ensign.rotational.io//my-org/hermes-international-sa/my-org")
  // })

  // it("returns ensign.rotational.io/baskin-robins when org name is Baskin Robins", () => {
  //     expect(slugify("Baskin-Robins", "my org")).toBe("ensign.rotational.io/baskin-robins/my-org")
  // })
});
