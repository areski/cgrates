ApierV1.SetTPRate
+++++++++++++++

Creates a new rate within a tariff plan.

**Request**:

 Data:
  ::

   type TPRate struct {
	TPid      string     // Tariff plan id
	RateId    string     // Rate id
	RateSlots []RateSlot // One or more RateSlots
   }

   type RateSlot struct {
	ConnectFee       float64 // ConnectFee applied once the call is answered
	Rate             float64 // Rate applied
	RatedUnits       int     //  Number of billing units this rate applies to
	RateIncrements   int     // This rate will apply in increments of duration
	GroupInterval    int     // Group position
	RoundingMethod   string  // Use this method to round the cost
	RoundingDecimals int     // Round the cost number of decimals
	Weight           float64 // Rate's priority when dealing with grouped rates
   }

 Mandatory parameters: ``[]string{"TPid", "RateId", "ConnectFee", "RateSlots"}``

 *JSON sample*:
  ::

   {
    "id": 1, 
    "method": "ApierV1.SetTPRate", 
    "params": [
        {
            "RateId": "SAMPLE_RATE_2", 
            "RateSlots": [
                {
                    "ConnectFee": 0.2, 
                    "Rate": 2, 
                    "RateIncrements": 60, 
                    "RatedUnits": 1, 
                    "RoundingDecimals": 2,
                    "GroupInterval": 0, 
                    "RoundingMethod": "*up", 
                    "Weight": 10.0
                }, 
                {
                    "ConnectFee": 0.2, 
                    "Rate": 2.1, 
                    "RateIncrements": 1, 
                    "RatedUnits": 1, 
                    "RoundingDecimals": 2,
                    "GroupInterval": 60, 
                    "RoundingMethod": "*up", 
                    "Weight": 20.0
                }
            ], 
            "TPid": "SAMPLE_TP"
        }
    ]
   }

**Reply**:

 Data:
  ::

   string

 Possible answers:
  ``OK`` - Success.

 *JSON sample*:
  ::

   {
    "error": null, 
    "id": 0, 
    "result": "OK"
   }

**Errors**:

 ``MANDATORY_IE_MISSING`` - Mandatory parameter missing from request.

 ``SERVER_ERROR`` - Server error occurred.

 ``DUPLICATE`` - The specified combination of TPid/RateId already exists in StorDb.


ApierV1.GetTPRate
+++++++++++++++

Queries specific rate on tariff plan.

**Request**:

 Data:
  ::

   type AttrGetTPRate struct {
	TPid   string // Tariff plan id
	RateId string // Rate id
   }

 Mandatory parameters: ``[]string{"TPid", "RateId"}``

 *JSON sample*:
  ::

   {
    "id": 1, 
    "method": "ApierV1.GetTPRate", 
    "params": [
        {
            "RateId": "SAMPLE_RATE_4", 
            "TPid": "SAMPLE_TP"
        }
    ]
   }
   
**Reply**:

 Data:
  ::

   type TPRate struct {
	TPid      string     // Tariff plan id
	RateId    string     // Rate id
	RateSlots []RateSlot // One or more RateSlots
   }

   type RateSlot struct {
	ConnectFee         float64 // ConnectFee applied once the call is answered
	Rate               float64 // Rate applied
	RateUnit           string     //  Number of billing units this rate applies to
	RateIncrement      string     // This rate will apply in increments of duration
	GroupIntervalStart string     // Group start time during a call
	RoundingMethod     string  // Use this method to round the cost
	RoundingDecimals   int     // Round the cost number of decimals
	Weight             float64 // Rate's priority when dealing with grouped rates
   }

 *JSON sample*:
  ::

   {
    "error": null, 
    "id": 2, 
    "result": {
        "RateId": "SAMPLE_RATE_2", 
        "RateSlots": [
            {
                "ConnectFee": 0.2, 
                "Rate": 2, 
                "RateIncrement": "60s", 
                "RateUnit": "1s", 
                "RoundingDecimals": 2,
                "GroupIntervalStart": "0s", 
                "RoundingMethod": "*up", 
                "Weight": 10
            }, 
            {
                "ConnectFee": 0.2, 
                "Rate": 2.1, 
                "RateIncrement": "1s", 
                "RateUnit": "1s", 
                "RoundingDecimals": 2,
                "GroupIntervalStart": "60s",
                "RoundingMethod": "*up", 
                "Weight": 20
            }
        ], 
        "TPid": "SAMPLE_TP"
    }
   }

**Errors**:

 ``MANDATORY_IE_MISSING`` - Mandatory parameter missing from request.

 ``SERVER_ERROR`` - Server error occurred.

 ``NOT_FOUND`` - Requested rate id not found.


ApierV1.GetTPRateIds
++++++++++++++++++

Queries rate identities on tariff plan.

**Request**:

 Data:
  ::

   type AttrGetTPRateIds struct {
	TPid string // Tariff plan id
   }

 Mandatory parameters: ``[]string{"TPid"}``

 *JSON sample*:
  ::

   {
    "id": 1, 
    "method": "ApierV1.GetTPRateIds", 
    "params": [
        {
            "TPid": "SAMPLE_TP"
        }
    ]
   }

**Reply**:

 Data:
  ::

   []string

 *JSON sample*:
  ::

   {
    "error": null, 
    "id": 1, 
    "result": [
        "SAMPLE_RATE_1", 
        "SAMPLE_RATE_2", 
        "SAMPLE_RATE_3", 
        "SAMPLE_RATE_4"
    ]
   }

**Errors**:

 ``MANDATORY_IE_MISSING`` - Mandatory parameter missing from request.

 ``SERVER_ERROR`` - Server error occurred.

 ``NOT_FOUND`` - Requested tariff plan not found.


