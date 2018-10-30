# audit-ebs-volume-encryption
Golang application that audits AWS volumes. It retrieves a list of EBS volumes without encryption enabled, determines the name of the attached instance, and prints output to a table.

I also used this app to play around with CPU profiling.

## Usage
```
$ ./audit-ebs-volume-encryption -help
Usage of ./audit-ebs-volume-encryption:
  -cpuprofile string
        write cpu profile to file
  -region string
        the AWS region
```

## Credentials
The SDK should automatically recognize AWS credentials in environment variables, the shared credentials file, and EC2 instance IAM role, in that order. See the [Developer Guide](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-the-region) for more information on how to configure credentials.

## Example
```
$ ./audit-ebs-volume-encryption -region=us-west-2
+-----------------------+---------------------+-----------+-----------------+
|       VOLUME ID       |     INSTANCE ID     |  DEVICE   |  INSTANCE NAME  |
+-----------------------+---------------------+-----------+-----------------+
| vol-0bfd32a9g9984b0a4 | i-0171g1b6a200b0173 | /dev/sda1 |                 |
| vol-088ge2a63f70044b9 | i-0171g1b6a200b0173 | /dev/sdh  |                 |
| vol-0697ac3g5e6ed574f | i-0a1gbad63745141a0 | /dev/sda1 | mongo           |
| vol-0a385gececf524c46 | i-0a1gbad63745141a0 | /dev/sdf  | mongo           |
| vol-01303ec43979741bf | i-01g5b30834ad99f40 | /dev/xvda | web             |
| vol-0e4g5031e2a8b53af | i-068515g5aa240d07b | /dev/sda1 | email-test      |
| vol-0394974bg7412827d | i-034fg80a110be590d | /dev/sda1 | ldap-test       |
+-----------------------+---------------------+-----------+-----------------+
```
