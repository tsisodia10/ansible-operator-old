#!/bin/bash

BASEDIR=$(dirname "$(readlink -f "$0")")

# Parse command line arguments
for i in "$@"
do
case $i in
   --user=*)
   USER=`echo $i | sed 's/[^env_var_*]'`
   ;;
   --sshkey=*)
   SSHKEY=`echo $i | sed 's/[-a-zA-Z0-9]*=//'`
   ;;
   --projectrepo=*)
   PROJECTREPO=`echo $i | sed 's/[-a-zA-Z0-9]*=//'`
   ;;
   --playbook=*)
   PLAYBOOK=`echo $i | sed 's/[-a-zA-Z0-9]*=//'`
   ;;
   --username=*)
   USERNAME=`echo $i | sed 's/[-a-zA-Z0-9]*=//'`
   ;;
   --inventory=*)
   INVENTORY=`echo $i | sed 's/[-a-zA-Z0-9]*=//'`
   ;;
   --password=*)
   PASSWORD=`echo $i | sed 's/[-a-zA-Z0-9]*=//'`
   ;;*)
   echo "Unknown option $i"        # unknown option
   HELP=True
   UNKNOWN=True
   ;;
esac
done

echo ""
echo "Container Setup"
echo "......................................................."
echo "BASEDIR                               ${BASEDIR}"
echo "USER                                  ${USER}"
echo "SSH_KEY                               ${SSH_KEY}"
echo "USERNAME                              ${USERNAME}"
echo "PASSWORD                              ${PASSWORD}"
echo "PLAYBOOK                              ${PLAYBOOK}"
echo "INVENTORY                             ${INVENTORY}"

${USERANME} > /env/vars


