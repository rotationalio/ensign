export const slugifyMockData = () => {
  return [
    {
      input: 'This is a test ---',
      expected: 'this-is-a-test',
    },
    {
      input: '___This is a test ---',
      expected: 'this-is-a-test',
    },
    {
      input: '___This is a test___',
      expected: 'this-is-a-test',
    },
    {
      input: 'This -- is a ## test ---',
      expected: 'this-is-a-test',
    },
    {
      input: '  THIS  is   a   test     ',
      expected: 'this-is-a-test',
    },
    {
      input: '影師嗎',
      expected: 'yingshima', // should be 'ying-shi-ma'
    },
    {
      input: "C'est déjà l'été.",
      expected: 'c-est-deja-l-ete',
    },
    {
      input: 'Nín hǎo. Wǒ shì zhōng guó rén',
      expected: 'nin-hao-wo-shi-zhong-guo-ren',
    },
    {
      input: 'jaja---lol-méméméoo--a',
      expected: 'jaja-lol-mememeoo-a',
    },
    {
      input: 'Компьютер',
      expected: 'kompyuter', // should be 'kompiuter'
    },
    {
      input: 'foo &amp; bar',
      expected: 'foo-bar',
    },
    {
      input: '10 amazing secrets',
      expected: '10-amazing-secrets',
    },
    {
      input: 'buildings with 1000 windows',
      expected: 'buildings-with-1000-windows',
    },
    {
      input: 'recipe number 3',
      expected: 'recipe-number-3',
    },
    {
      input: '404',
      expected: '404',
    },
    {
      input: '1,000 reasons you are #1',
      expected: '1000-reasons-you-are-1',
    },
    {
      input: 'I ♥ 🦄',
      expected: 'i',
    },
    {
      input: 'i love 🦄',
      expected: 'i-love',
    },
  ];
};
