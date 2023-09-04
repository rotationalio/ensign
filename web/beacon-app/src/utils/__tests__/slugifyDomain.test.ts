/* eslint-disable prettier/prettier */
import { describe, expect, it } from 'vitest';

import { slugifyMockData } from '../__mocks__';
import { slugify, stringify_org } from '../slugifyDomain';
describe('#slugify', () => {
  it('returns https://rotational.app when org name is empty', () => {
    expect(slugify('')).toBe('https://rotational.app');
  });

  it('returns https://rotational.app/rotational-labs-inc/domain-space when org name is Roational Labs, Inc.', () => {
    expect(slugify('domain space', 'Rotational Labs, Inc.')).toBe(
      'https://rotational.app/rotational-labs-inc/domain-space'
    );
  });

  slugifyMockData().map(({ input, expected }) => {
    it(`should slugify ${input} to ${expected}`, () => {
      expect(stringify_org(input)).toBe(expected);
    });
  });
});
