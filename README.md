# pythonx
Let python be easier to use.

## Reason
I made a speech in PyCon China 2021. I have raised an idea of 'pythonx'. But I haven't implement it yet.

Before I get started, let me just do a planning.

## Goal
* It's a command line tool
* It can detect '#python3.9' tag on the top of a .py file; If that version of Python does not exist in local, we download and install it; If that version of Python exist in local, we use it to run the .py file
* It can detect 'requirements.txt' file on the same folder level of .py file you want to run; If thoese dependencies does not exist in the corresponse version of Python, we install it

## Plan

* We use Golang to make the binary executable file because it's natural (And we don't want to use Python to install Python, because python's version is different across different machines and systems.
* It has permission to edit '/home', becuase that's where we put '.pythonx' folder into; Inside the '.pythonx' folder, we manage different python version for our users. Sprated by version tag, like 'python3.9', 'python3.6'
* I even doubt for such a simple need, do we really need this package? https://github.com/urfave/cli/blob/master/docs/v2/manual.md
* We'll use ; We'll use 'gox' for the cross-compiling
