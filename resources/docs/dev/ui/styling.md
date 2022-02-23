# Styling
For styling this project uses the preprocessor scripting language [Sass](https://sass-lang.com/). It provides more utility than native CSS by allowing the use of variables, functions, imports, and more.

## Styling Convention
Each component that has styling should have it's own `.scss` file in its directory. The naming of classes for elements should follow the convention of `ComponentName` for the returned component, `ComponentName__section` for sub-elements, and `ComponentName__section--modifier` for modifiers. I.e. `Dataset__results--warning`

## Hoss Theme Support
The Hoss is setup to allow for a custom theme using a primary and secondary color. Due to the nature of Sass as a preprocessor, this theme has to be defined using native CSS variables. This is currently being set as `--main-color` and `--secondary-color` upon response from the `ui/config.json` endpoint, which is fetched on load.

When adding new SVG assets you must modify them to use `var(--main-color)` in order to match the applications theme.  Due to limitations with CSS, the SVG's cannot be assigned as backgrounds as the CSS variables will not work and will therefore render incorrectly.  To correctly add a SVG it must be imported as a React component into the TSX.
