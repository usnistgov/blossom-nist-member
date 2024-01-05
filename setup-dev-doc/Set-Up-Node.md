## The dashboard projects utilize Node/TypeScript and React/Vite frameworks.

### The method of Node deployment **preferred by the team** uses NVM for Node version management

If you want to have easier way to manage your Node versions you can install NVM in the following way: 

### For Brew-based installation on MacOS

1. Install NVM 
```
brew install nvm
``` 

2. Create NVM working directory (usually ~/.nvm, which is the same as $HOME/.nvm)
```
mkdir ~/.nvm
```

3. Append the NVM to **your default shell** (the Latest MacOS default is zsh, though yours can be another, e.g. bash). To determine your default shell inside your terminal application (iTerm/iTerm2, bash, zsh, etc.) type ``` echo $SHELL ``` if response is ``` /bin/bash ``` then your shell is bash, if your response is ``` /bin/zsh ```, then your shell is zsh. 

-------
### For reference and to avoid confusion Bash (and partially ZSH) Startup Order is described below

#### When Bash is started for **OS login in terminal or launched with -l/--login flag** (and a few other rare contexts) Bash will read configuration files in the following order: 

   1. ```/etc/profile``` This file (applies to **configuration for ALL users**) 
   2. ```~/.bash_profile```  ✕ (if this file found => the items below would be ignored) ✕
   3. ```~/.bash_login```  ✕ (if this file found => the items below would be ignored)  ✕
   4. ```~/.profile```  ✕ (if this file found => the items below would be ignored)  ✕
   
The above behavior can be stopped by flag --noprofile. To find out more read applicable sections in  ```man bash ```

#### When Bash is started as **authenticated user's interactive session** the home directory file read by bash is ```~/.bashrc``` (for zsh it is ```~/.zshrc```), which alternatively can be force-read by using flag ```--rcflie```


[This unix.stackexchange.com answer](https://unix.stackexchange.com/questions/439042/debian-read-order-of-bash-session-configuration-files-inconsistent) describes how one can include optionally read profile-files matryoshka (AKA  russian-nested-doll) style to avoid default behaviour pit-falls and make your Bash behaviour predictably consistent.

-------

a. 
```

```

b. 
```

```

-----------

You should create NVM's working directory if it doesn't exist:
mkdir ~/.nvm

Add the following to your shell profile e.g. ~/.profile or ~/.zshrc:
  export NVM_DIR="$HOME/.nvm"
  [ -s "/opt/homebrew/opt/nvm/nvm.sh" ] && \. "/opt/homebrew/opt/nvm/nvm.sh"  # This loads nvm
  [ -s "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm" ] && \. "/opt/homebrew/opt/nvm/etc/bash_completion.d/nvm"  # This loads nvm bash_completion

You can set $NVM_DIR to any location, but leaving it unchanged from
/opt/homebrew/Cellar/nvm/0.39.7 will destroy any nvm-installed Node installations
upon upgrade/reinstall.
