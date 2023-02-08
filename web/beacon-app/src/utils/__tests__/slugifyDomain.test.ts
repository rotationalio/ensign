/* eslint-disable prettier/prettier */
import { describe, expect, it } from 'vitest';

import { slugify } from '../slugifyDomain';

describe('#slugify', () => {
  it('returns https://rotational.app when org name is empty', () => {
    expect(slugify('')).toBe('https://rotational.app');
  });

  it('returns https://rotational.app/rotational-labs-inc/domain-space when org name is Roational Labs, Inc.', () => {
    expect(slugify('domain space', 'Rotational Labs, Inc.')).toBe(
      'https://rotational.app/rotational-labs-inc/domain-space'
    );
  });

  // it("returns https://rotational.app/hermes-international-sa when org name is Hermès International S.A.", () => {
  //     expect(slugify("Hermès International S.A.", "my org")).toBe("rotational.app//my-org/hermes-international-sa/my-org")
  // })

  // it("returns https://rotational.app/baskin-robins when org name is Baskin Robins", () => {
  //     expect(slugify("Baskin-Robins", "my org")).toBe("rotational.app/baskin-robins/my-org")
  // })
});
