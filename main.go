// Init the login service
// Init the graph traversal queue
// Add the first page to the queue
// Run the main loop with wait group
// Print result

// Main loop:
// Get login
// Get page
// Register response
// Enqueue children
// Dequeue anything that is in the queue and spawn new goroutines (increment wg)
// wg.Done