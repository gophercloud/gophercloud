/*
Package webhooks provides the ability to trigger an action represented by a webhook from the OpenStack Clustering
Service.

Example to Trigger webhook action

	result, err := webhooks.Trigger(webhookClient, "f93f83f6-762b-41b6-b757-80507834d394", nil).Extract()
	if err != nil {
		panic(err)
	}
*/
package webhooks
