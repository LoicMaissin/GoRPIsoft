package analyser

/*
  Functions to translate the actuator's response (byte array)
  to human-readable data
*/

func nthBit(x byte, n uint) bool {
	if x>>n&1 == byte(0) {
		return false
	}
	return true
}

func addBytes(v []byte) int {
	// The number equal to v concatenated, with v[0] MSB
	res := 0
	for _, x := range v {
		res = res<<8 + int(x)
	}
	return res
}

// ActuatorInfo Translate the actuator's response to human-readable data
func ActuatorInfo(response [38]byte) map[string]interface{} {

	fields := map[string]interface{}{
		"isOpened":                          nthBit(response[5], 0),
		"isClosed":                          nthBit(response[5], 1),
		"torqueLimiterActionOpenDirection":  nthBit(response[5], 2),
		"torqueLimiterActionCloseDierction": nthBit(response[5], 3),
		"selectorToLocalPosition":           nthBit(response[5], 4),
		"selectorToRemotePosition":          nthBit(response[5], 5),
		"selectorToOffPosition":             nthBit(response[5], 6),
		"powerOn":                           nthBit(response[5], 7),

		"actOpening":               nthBit(response[6], 0),
		"actClosing":               nthBit(response[6], 1),
		"handwheelAction":          nthBit(response[6], 2),
		"ESDCommand":               nthBit(response[6], 3),
		"actRunning":               nthBit(response[6], 4),
		"actFault":                 nthBit(response[6], 5),
		"positionSensorPowerFault": nthBit(response[6], 6),
		"torqueSensorPowerFault":   nthBit(response[6], 7),

		"lockedMotorOpen":      nthBit(response[7], 0),
		"lockedMotorClose":     nthBit(response[7], 1),
		"motorThermalOverload": nthBit(response[7], 2),
		"lostPhase":            nthBit(response[7], 3),
		"overtravelAlarm":      nthBit(response[7], 4),
		"directionOpenAlarm":   nthBit(response[7], 5),
		"directionCloseAlarm":  nthBit(response[7], 6),
		"batteryLow":           nthBit(response[7], 7),

		"runningTorque": response[8],
		"actPosition":   response[9],

		"indication1":       nthBit(response[10], 0),
		"indication2":       nthBit(response[10], 1),
		"indication3":       nthBit(response[10], 2),
		"indication4":       nthBit(response[10], 3),
		"indication5":       nthBit(response[10], 4),
		"valveJammed":       nthBit(response[10], 5),
		"auxiliary24VFault": nthBit(response[10], 6),
		"tooManyStarts":     nthBit(response[10], 7),

		"pumping":              nthBit(response[11], 0),
		"confMemFault":         nthBit(response[11], 1),
		"activityMemFault":     nthBit(response[11], 2),
		"baseMemFault":         nthBit(response[11], 3),
		"stopMidTravel":        nthBit(response[11], 4),
		"lostSignal":           nthBit(response[11], 5),
		"partialStrokeRunning": nthBit(response[11], 6),
		"partialStrokeFault":   nthBit(response[11], 7),

		"openBreakoutMaxTorque": response[12],
		"closeTightMaxTorque":   response[13],
		"openingMaxTorque":      response[14],
		"closingMaxTorque":      response[15],
		"startsLast12h":         addBytes(response[16:18]),
		"totalStarts":           addBytes(response[18:22]),
		"totalRunningTime":      addBytes(response[22:26]),
		"partialStarts":         addBytes(response[26:30]),
		"partialRunningTime":    addBytes(response[30:34]),
		"actPosition(per mil)":  addBytes(response[34:36]),
	}
	return fields
}
