/*
Package webhooks enables triggers an action represented by a webhook from the OpenStack
Clustering Service.

Example to Trigger webhook action
	result := webhooks.Trigger(webhookClient, "my-webhook")
	if result.Err != nil {
		panic(result.Err)
	}
*/
package webhooks
