# EC2 Connector

A command-line tool that simplifies connecting to Amazon EC2 instances via SSH. It utilizes [fzf](https://github.com/junegunn/fzf) to help you select the instance you want to connect to and offers several configurable options, including the SSH user, the region to search for instances, and the directory where your SSH keys are stored.

### Prerequisites

- Install [FZF](https://github.com/junegunn/fzf): Make sure you have ```fzf``` installed on your system.
- Set Up AWS Credentials: Follow the instructions to set up your AWS credentials [here](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html).

### Build from Source
  
To install the latest version of ec2-connect, you can either download the precompiled [binaries](https://github.com/triplemcoder14/ec2-connect/releases/tag/v1.2.1)  or build it from source: run

```
 go install github.com/triplemcoder14/ec2-connect@latest
```

### Usage

#### Running the Tool
You can run the tool directly using:

```
ec2-connect
```

### Command-Line Options

You can specify options when running the command:

- ``user``: SSH user for login (default is ``ubuntu``).
- ``directory``: Directory where your SSH keys are stored (default is ``~/.ssh/id_rsa``).
- ``region``: AWS region to use (default is ``us-east-2``).




```
ec2-connect [OPTIONS]
```
### Options

- ``user``: The SSH user for login. The default is ``ubuntu``.
- ``directory``: The directory where your SSH keys are stored. Default is ```~/.ssh```.
- ``region``: The AWS region where your EC2 instances are located. Default is ``us-east-2``.

  Example:

```
  ec2-connect -user ec2-user -directory /home/ubuntu/keys -region us-east-2
```

### Interactive Prompts

If you run the tool without options, it will prompt you to:

- Select an Instance: Enter the Instance ID of the EC2 instance you want to connect to.
- Select a Region: If no region is specified, it will show a list of available regions via fzf.

### Connecting to an EC2 Instance
- Ensure Your Instance is Running: Check the AWS Management Console to verify that your EC2 instance is in a running state and has a public IP.
- Run the Tool: Use the command with or without flags as described above.

#### Troubleshooting

- ``No Running Instances Found``: Ensure that your AWS account has running instances in the selected region.
- ``SSH Key Issues``: Make sure your SSH private key is in the specified directory and has the correct permissions (e.g., chmod 400 ``/path/to/your/key.pem``)


### Example Commands

- Basic Usage:

  
```
ec2-connect
```

- With User and Directory:
  
``
ec2-connect -user ec2-user -directory /home/ubuntu/keys
``

- With Region Specification:
  
```
ec2-connect -region us-east-2
```



### License

This tool is released under the Apache License. See the [LICENSE](https://github.com/triplemcoder14/ec2-connect/blob/main/LICENSE) file for more information.

