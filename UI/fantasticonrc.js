module.exports = {
  inputDir: 'src/assets/icons', // (required)
  outputDir: 'src/assets/font-icon', // (required)
  fontTypes: [],
  assetTypes: ['ts', 'css', 'json', 'html'],
  fontsUrl: '/static/fonts',
  formatOptions: {
    // Pass options directly to `svgicons2svgfont`
    woff: {
      // Woff Extended Metadata Block - see https://www.w3.org/TR/WOFF/#Metadata
      metadata: '...'
    },
    json: {
      // render the JSON human readable with two spaces indentation (default is none, so minified)
      indent: 2
    },
    ts: {
      // select what kind of types you want to generate (default `['enum', 'constant', 'literalId', 'literalKey']`)
      types: ['constant', 'literalId'],
      // render the types with `'` instead of `"` (default is `"`)
      singleQuotes: true
    }
  }
};
