# Premises/requisites/considerations
*	Agents don't send updates frequently so they can wait up to 5sec (read ≫ write)
*	It's better to put clusters in a another collection than ```hosts``` because they are used in licensing logic 
# Design
* ```hosts``` collection
	* Contains all hostdata
* ```alerts``` collection
* ```clusters``` collection