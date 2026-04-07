import sveltePreprocess from 'svelte-preprocess';

const contentReplacer = args => {
  return global.zPlaceholderReplacer && { code: global.zPlaceholderReplacer(args.content, global.zIsSSR) };
};

export default {
  preprocess: [
    sveltePreprocess(),
    {
      style: contentReplacer,
      script: contentReplacer,
      markup: contentReplacer,
    },
  ],
  cache: false,
  compilerOptions: { preserveComments: true },
};
