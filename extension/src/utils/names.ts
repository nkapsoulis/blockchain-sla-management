function kebabCase(word: string) {
  return word
    .toLowerCase()
    .replace(':', '')
    .replace(/\s+/g, '-');
}

const names = {
  kebabCase,
};

export default names;
