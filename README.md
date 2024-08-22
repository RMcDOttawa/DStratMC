<h1>DStratMC - Dart Strategy Monte-Carlo</h1>
<h2>Introduction</h2>
This program began as a learning exercise - it seemed it would be fun, and it was.
<p>When playing darts, certain areas on the board are well-defined as the highest scoring areas. 
E.g. the "treble 20" segment is the highest score on the board, then centre "red bull" is 
the next highest, etc.
<p>However, this doesn't necessarily mean that you should aim at those points. 
Every player has a certain degree of accuracy in their throw.  
<ul><li>Good players, especially professions, generally hit what they are aiming at, s
o they will certainly aim for treble-20 during the game phase when they are trying to maximize their scoring.</li>
<li>Beginners, however, don't always hit what they aiming for.  
A dart aimed at a particular point on the board is likely to fall in a circle, 
centered on that point.  The more the beginner, the larger the circle of that likely hit.</li>
</ul>
<p>The dartboard is deviously designed with this in mind.  
The <i>lowest</i> scoring segment on the board is <i>right next door</i> to the highest segment. 
So a beginner aiming for treble-20 is likely to land in the "1" segment quite a bit of the time or, on the other side of the 20, the "5" segment.
<p>Classic advice to beginners is to aim for treble-19 instead, because the other segments 
that the darts might fall into are higher scoring.
<p>This seemed like a good topic for experimentation.
<ul>
<li>If players "circles of accuracy" get smaller as they improve, then there must be a 
threshold point where aiming for treble-20 or treble-19 is, statistically, the smarter move.</li>
<li>And what about other options?  Depending on the size of one's accuracy circle, are 
there other good areas on the board to aim at?</li>
</ul>
<p>This seemed like a good use of a Monte Carlo technique: a program that simulates 
throwing a large number of darts at a target and, using a model that simulates a player's accuracy, calculates an average score.  Then, the program can throw a large number of darts at each of a large number of targets, and find the high-scoring areas on the board.
<h2>Instructions</h2>
Run the program in your GO IDE - or use the compiled binary if there is one available 
for your platform.  To run from the IDE, run the main/main.go file.
<p>You will see a dartboard on the right and a set of controls on the left.
<p>The "Interaction Mode" radio buttons control the overall behaviour of the program - 
especially what happens when you click the mouse on the board.
<table>
<thead>
<tr>
<td width="200px"><b>Mode</b></td>
<td><b>Behaviour</b></td>
</thead>
<tbody>
	<tr valign="top">
		<td>One Throw Exact</td>
		<td>When you click on the board, you get feedback of the dart landing exactly
			where you aimed it, and the resulting score. This is not useful for anything except 
			verifying that the UI and interaction are working.</td>
	</tr>
	<tr valign="top">
		<td width="200px">One Throw Normal</td>
		<td>Simulates throwing a dart that falls in an area determined by your degree of accuracy. 
		The possible landing spots are centered on the target where you click, and the amount 
		of variation from that spot is determined by your "standard deviation", which you
		enter in the field so labelled on the screen.
		<p>Standard Deviation is measured as a fraction between 0 and 1 and refers, roughly,
		to a circle that is that fraction of the diameter of the board. So, for example, 
		a standard deviation of 0.25 means that about 95% of your darts will land within
		a circle that is about 1/4 of the width of the board.
		<p>The mathematical model is called a <i>Normal Distribution</i>, which is what the
		word "Normal" is referring to in the interaction modes.
		</td>
	</tr>
	<tr valign="top">
		<td width="200px">Multi Throw Normal</td>
		<td>In this mode, when you click on the board, a large number of darts are thrown
		at the target, landing with the same normal distribution.  The number of darts
		thrown can be entered in the field labeled "throws".
		With the dots drawn to indicate where your throws land, you will be able to see
		the normal distribution - the darts land in a circle around your target, with 
		more toward the centre.</td>
	</tr>
	<tr valign="top">
		<td width="200px">Search Normal</td>
		<td>Finally, with this setting, you don't click on the board. Instead, just click
		on the "search" button. The program will throw a large number of darts at many locations
		around the board, and will report back to you on the location of the 10 best throws.</td>
	</tr>
</tbody>
</table>
<p>Note that some other interaction modes are implemented in the code, but are commented out. 
Feel free to turn them back on for experimentation.
<p>The 3 checkboxes labeled "show circles for..." will display 3 circles on the board, showing
where 1, 2, and 3 standard deviations lie from your target point. According to the statistics,
about 68% of your throws will land within the 1 standard deviation circle, about 95% within
the 2 standard deviation circle, and about 99.7% within the 3 standard deviation circle.
With a large number of throws, a few will even fall outside that circle - wild darts do happen.
<h2>Future Feature</h2>
<p>Right now, you have to enter the Standard Deviation figure manually.
<ul>
<li>Less than 0.1 is a darn good player.</li>
<li>0.1 is a good player like the guy at the pub that you can't beat.</li>
<li>0.2 is a beginner with some experience, starting to hit what they aim for fairly often.</li>
<li>0.3 is a true beginner who doesn't really have that much control over where their darts land, 
and misses the board entirely from time to time</li>
</ul>
<p>Future releases of this program will:
<ul>
<li>Allow you to <i>draw</i> your circle of accuracy on the board instead of having to
type in a number, and </li>
<li>Measure your own standard deviation by telling the program what target you are aiming
for, then clicking where your dart actually lands. Throwing 20 or so darts this way will give
a good working figure of your own personal standard deviation.</li>
</ul>