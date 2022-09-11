<h3>Game Algorithm</h3>

<ol><li><p>Initialize Screen</p></li><li><p>Initialize Game Objects</p><ol><li><p>Set initial values for snake variable</p></li><li><p>Set initial values for apple variable</p></li></ol></li><li><p>Display Game Frame</p><ol><li><p">Based on frame’s width and height call function to draw game’s border</p></li></ol></li><li><p>Display Initial Game Score </p></li><li><p>Run background function to continuously listen for user input</p></li><li><p>If <code>isGameOver</code> is set to true then to step 13</p></li><li><p>If <code>isGamePaused</code> set to true then</p><ol><li><p>Display Game Pause Info</p></li></ol></li><li><p>Read user input and perform accordingly</p><ol><li><p>If no user input go to step 9</p></li><li><p>If user input is <code>q</code> exit the game</p></li><li><p>If user input is <code>p</code> toggle value of <code>isGamePaused</code></p></li><li><p>If user input if up arrow key and snake is moving horizontally then set snake’s <code>rowVelocity</code> to -1 and <code>columnVelocity</code> to 0</p></li><li><p>If user input if down arrow key and snake is moving horizontally then set snake’s <code>rowVelocity</code> to 1 and <code>columnVelocity</code> to 0</p></li><li><p>If user input if left arrow key and snake is moving vertically then set snake’s <code>rowVelocity</code> to 0 and <code>columnVelocity</code> to -1</p></li><li><p>If user input if right key and snake is moving vertically then set snake’s <code>rowVelocity</code> to 0 and <code>columnVelocity</code> to 1</p></li></ol></li><li><p>Update game state</p><ol><li><p>If <code>isGamePaused</code> set to true don’t do anything; return</p></li><li><p>Clear the screen</p><ol><li><p>Based on <code>coordinatesToClear</code> data call print function to display empty space on desired coordinates</p></li></ol></li><li><p>Update snake variable’s values</p><ol><li><p>Get snake’s current head coordinates</p></li><li><p>Create new coordinate by adding snake head’s x-coordinate by snake’s <code>columnVelocity</code> and adding snake head’s y-coordinate by snake’s <code>rowVelocity</code></p></li><li><p>Add new coordinate to snake’s point field</p></li><li><p>Set snake’s coordinates within game frame</p><ol><li><p>Get game frame’s top left x and y coordinate</p></li><li><p>Determine game frame’s boundaries as follows</p><ol><li><p>left boundary is same as frame’s top left x coordinate</p></li><li><p>top boundary axis is same as frame’s top left y coordinate</p></li><li><p>right boundary is equal to sum of left boundary and frame’s with minus frame’s boundary thickness</p></li><li><p>bottom boundary is equal to sum of top boundary and frame’s height</p></li></ol></li><li><p>For each snake’s coordinate update snake’s coordinate such that it is inside game’s frame</p><ol><li><p>If snake’s y coordinate is less than or equal to top boundary then set new y coordinate as bottom boundary - 1</p></li><li><p>If snake’s y coordinate is greater than or equal to bottom boundary then set new y coordinate as top boundary + 1</p></li><li><p>If snake’s x coordinate is less than or equal to left boundary then set new x coordinate as right boundary - 1</p></li><li><p>If snake’s x coordinate is greater than or equal to right boundary then set new x coordinate as left boundary + 1 </p></li></ol></li></ol></li><li><p>Check if snake ate apple</p><ol ><li><p>For each coordinate of snake check if it is same as apple’s coordinate</p><ol><li><p>If any match is found then </p><ol><li><p>increase <code>score</code> by 1</p></li><li><p>Call function that updates rendering for game score</p></li></ol></li><li><p>Else</p><ol><li><p>Append snake’s first coordinate to <code>coordinatesToClear</code> variable</p></li><li><p>Slice snake’s <code>points</code> to range from index 1 to end</p></li></ol></li></ol></li></ol></li><li><p>Check if snake is eating itself</p><ol><li><p>Get snake’s head coordinates</p></li><li><p>For all other coordinates of snake check if any is equal to snake’s head coordinates</p></li><li><p>If any match is found then set <code>isGameOver</code> to true</p></li></ol></li></ol></li><li><p>Update apple variable’s value</p><ol><li><p>As long as for all coordinates of snake if any matches with apple’s coordinate generate a new coordinate for apple</p></li></ol></li></ol></li><li><p>Display game objects</p><ol><li><p>Display all coordinates of snake in screen</p></li><li><p>Display coordinate of apple in screen</p></li></ol></li><li><p>Wait for 75ms so that game isn’t too fast for human eye</p></li><li><p>Go to step 6</p></li><li><p>Display game over info</p></li><li><p>Exit game</p></li></ol>