# ec2-rolling-snapshot


Make sure to configure your AWS credentials:
https://github.com/aws/aws-sdk-go#configuring-credentials


Quick Notes:

 * Relys on the snapshot-task name for ensuring misc snapshots of the same volume aren't deleted.