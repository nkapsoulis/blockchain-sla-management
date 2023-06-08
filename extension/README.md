# Hyperledger Fabric Extension

Chrome extension written in React.

## Available Scripts

In the project directory, you can run:

### `npm run develop`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

Keep in mind that given that this is a chrome extension, which uses chrome specific APIs,
some things will not work. This function is still useful for examining the views.

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
The extension is ready to be deployed!

## Using the Extension

In order to test the app as an extension you can load it into a Chromium-based browser (like Chrome or Edge).
To do that you need to enable developer mode for extensions and then load the unpacked extension:

1. Go to Extensions page on your browser.
2. Use the slider to enable developer mode.
3. Click the load unpacked option and load the extension.
4. Try it out.
