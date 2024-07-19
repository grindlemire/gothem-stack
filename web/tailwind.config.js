const {
    iconsPlugin,
    getIconCollections,
} = require("@egoist/tailwindcss-icons");

/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["**/*.templ"],
    darkMode: "class",
    theme: {
        extend: {
            fontFamily: {
                mono: ["Courier Prime", "monospace"],
            },
        },
    },
    plugins: [
        require("daisyui"),
        iconsPlugin({
            collections: getIconCollections(["ic", "mdi"]),
        }),
    ],
    daisyui: {},
};
