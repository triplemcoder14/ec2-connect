# EC2 Connector

A cmd tool that simplifies connecting to Amazon EC2 instances via SSH. It utilizes [fzf](https://github.com/junegunn/fzf) to help you select the instance you want to connect to and offers several configurable options, including the SSH user, the region to search for instances, and the directory where your SSH keys are stored.

### Prerequisites

- Install FZF: Make sure you have ```fzf``` installed on your system.
- Set Up AWS Credentials: Follow the instructions to set up your AWS credentials [here](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html).

  ### Build from Source
To install the latest version, run:

```
 go install github.com/triplemcoder14/ec2-connect@latest
```

### Usage

To use the EC2 Connect, simply run the ec2-connect command and follow the prompts to select the EC2 instance you wish to connect to. By default, the tool uses the ubuntu user and looks for SSH keys in the ```~/.ssh directory.```

### Command Syntax

```
ec2-ssh [OPTIONS]
```
### Options

- ``user``: The SSH user for login. The default is ``ubuntu``.
- ``directory``: The directory where your SSH keys are stored. Default is ```~/.ssh```.
- ``region``: The AWS region where your EC2 instances are located. Default is ``us-east-1``.

  ### Example Command

  To connect to an EC2 instance using the ``ec2-user``, SSH keys stored in ``/home/ubuntu/keys``, and targeting the ``us-west-2 region``, you can use:

```
  ec2-ssh -user ec2-user -directory /home/ubuntu/keys -region us-west-2
```

### License

This tool is released under the MIT License. See the [LICENSE](https://github.com/triplemcoder14/ec2-connect/blob/main/LICENSE) file for more information.

