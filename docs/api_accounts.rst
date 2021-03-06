ApierV1.SetAccountActions
+++++++++++++++++++++++++


Process dependencies and load a specific AccountActions profile from storDb into dataDb.

**Request**:

 Data:
  ::

   type AttrSetAccountActions struct {
	TPid            string
	AccountActionsId string
   }

 Mandatory parameters: ``[]string{"TPid", "AccountActionsId"}``

 *JSON sample*:
  ::

   {
    "id": 0, 
    "method": "ApierV1.SetAccountActions", 
    "params": [
        {
            "AccountActionsId": "SAMPLE_AA_1", 
            "TPid": "SAMPLE_TP_1"
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
