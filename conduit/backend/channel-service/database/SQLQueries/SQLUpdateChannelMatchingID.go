package SQLQueries

// SQLChannelUpdateMatchingID updates channel row with matching ID
// Will update with channel fields
// channel_name
// channel_description
// status
func SQLChannelUpdateMatchingID() string {
	// not positive if I need to send the error, but hey why not
	return sqlChannelUpdateMatchingID
}

//TODO figure out what your going to with the status field and
const sqlChannelUpdateMatchingID = `	
UPDATE channels SET
    channel_name = ($2),
    channel_description = ($3),
    status = ($4),
    updated_at = now()
WHERE
    channel_id::text = ($1)             
AND
    deleted_at IS NULL
`

//NOTE: I have not put in the type in this because that will cause issues later, I'd rather just have the
// whole channel discarded cause transferring them to a whole new one is why to different
