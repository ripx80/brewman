package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ripx80/recipe"
	"github.com/ripx80/recipe/pkgs/m3w"
	"github.com/spf13/cobra"
)

// recipeCmd represents the recipe command
var recipeCmd = &cobra.Command{
	Use:   "recipe",
	Short: "handle and convert your brew recipes",
	Long:  `without any subcommand recipe convert m3 recipes to internal format`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		r, err := recipe.LoadFile(args[0], &recipe.M3{})
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(r)
		//r.SavePrettyYaml(args[1])
		// if *format == "yaml" {
		// 	err = recipe.SavePrettyYaml(*out)

		// } else {
		// 	err = recipe.SavePretty(*out)
		// }
	},
}

var recipeScale = &cobra.Command{
	Use:   "scale",
	Short: "recipe scale",
	Long: `scale the given recipe to water and sudyield
	usage: <recipefile> <water> <sudyield>
	`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		water, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		yield, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		rec, err := recipe.LoadFile(args[0], &recipe.Recipe{})
		if err != nil {
			fmt.Println(err)
			return
		}
		rScale, err := rec.Scale(water, yield)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(rScale)
	},
}

var recipeM3Down = &cobra.Command{
	Use:   "down",
	Short: "recipe down",
	Long:  `download all recipes from m3`,
	Run: func(cmd *cobra.Command, args []string) {
		outdir := "recipes"
		if len(args) != 0 {
			outdir = args[0]
		}
		m3w.Down("https://www.maischemalzundmehr.de", outdir)
	},
}

func init() {
	recipeCmd.AddCommand(recipeScale)
	recipeCmd.AddCommand(recipeM3Down)
}

/*cmd
recipe scale <water:27> <yield:65>
recipe yield <water:27> <originalwort:12> <filling:5.4kg>
Stammwürze ändern?
*/
