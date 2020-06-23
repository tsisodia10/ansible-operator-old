set -e

#Check extravars and print 
[[ ! -f /runner/env/extravars ]] && echo "---" > /runner/env/extravars
env | grep '^EXTRAVAR_' | while read EV ; do
    EV_KEY=$(echo "$EV" | awk -F '=' '{ print $1; }' | sed 's/^EXTRAVAR_//'
    EV_VALUE=$(echo "$EV" | awk -F '=' '{ print $2; }'
    echo "${EV_KEY}: ${EV_VALUE}" >> /runner/env/extravars
done

[[ ! -f /runner/env/ennvars ]] && echo "---" > /runner/env/envvars
env | grep '^ENVVAR_' | while read ENV ; do
    ENV_KEY=$(echo "$ENV" | awk -F '=' '{ print $1; }' | sed 's/^ENVVAR_//'
    ENV_VALUE=$(echo "$EV" | awk -F '=' '{ print $2; }'
    echo "${ENV_KEY}: ${ENV_VALUE}" >> /runner/env/extravars
done

[[ ! -f /runner/env/extravars ]] && echo "---" > /runner/env/extravars
env | grep '^EXTRAVAR_' | while read EV ; do
    EV_KEY=$(echo "$EV" | awk -F '=' '{ print $1; }' | sed 's/^E_//'
    EV_VALUE=$(echo "$EV" | awk -F '=' '{ print $2; }'
    echo "${EV_KEY}: ${EV_VALUE}" >> /runner/env/extravars
done


if [! -f /runner/inventory/hosts ];
     then echo "ERROR: Inventory is not mounted"
     exit 1;

if [! -f /runner/env/passwords && -f /runner/env/ssh_key ];
     then exit 1;

echo "Cloning GIT repository..............."

if [! -z {{PROJECT_TYPE}}];
      then echo "ERROR: Project type not specified"
      exit 1;

if [! -z {{PROJECT_URI}}];
      then echo "ERROR: Project URI not specified"
      exit 1;
      
if [! ! -d .git ] ;
      then echo "ERROR: This isn't a git directory" && return false;
fi

git_url=`git config --get remote.origin.url`
  if [[ $git_url != git@github.com:* && $git_url != git@gitlab.com:* ]]; then
    echo "ERROR: Remote origin is invalid" && return false;
fi

url=${git_url%.git};

if [[ $url == git@github.com:* ]]; then
    url=$(echo $url | sed 's,git@github.com:,https://github.com/,g');
  elif [[ $url == git@gitlab.com:* ]]; then
    url=$(echo $url | sed 's,git@gitlab.com:,https://gitlab.com/,g');
fi
open $url;

git clone $url; > /runner/project
echo "Git repository cloned"

if [ -f /runner/project/runner_playbook ];
    then echo "ERROR: No playbooks"
	    exit 1;

exec ansible-runner -- "{$@}"	    
