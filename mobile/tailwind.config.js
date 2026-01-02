/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./app/**/*.{js,jsx,ts,tsx}", "./components/**/*.{js,jsx,ts,tsx}"],
    theme: {
        extend: {
            colors: {
                primary: '#4ADE80',
                secondary: '#D97757',
                accent: '#F59E0B',
                'surface-light': '#F3F4F6',
                'surface-dark': '#1F2937',
                'card-light': '#FFFFFF',
                'card-dark': '#1F2937',
            },
            borderRadius: {
                'xl': '16px',
                '2xl': '24px',
            }
        },
    },
    plugins: [],
}
