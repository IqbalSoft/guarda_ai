package main

func generateBar(percentage int) string {
	filledBlocks := percentage / 10
	if filledBlocks > 10 {
		filledBlocks = 10
	}
	if filledBlocks < 0 {
		filledBlocks = 0
	}
	emptyBlocks := 10 - filledBlocks

	bar := ""
	for i := 0; i < filledBlocks; i++ {
		bar += "█"
	}
	for i := 0; i < emptyBlocks; i++ {
		bar += "░"
	}
	return bar
}
