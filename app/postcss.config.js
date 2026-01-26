import autoprefixer from 'autoprefixer'
import tailwindcss from "@tailwindcss/postcss"; /** Add this  */

export default {
  plugins: [
    // https://github.com/postcss/autoprefixer
    autoprefixer({
      overrideBrowserslist: [
        'last 4 Chrome versions',
        
      ]
    }),
    tailwindcss(), // ← Add this as well

  ]
}