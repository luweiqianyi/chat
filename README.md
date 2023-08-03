# chat
build connection between local program and remote git repository, commands are below. 
```shell
echo "# chat" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin https://github.com/luweiqianyi/chat.git
git push -u origin main
```


## reference
* `git`提交规范: `https://zhuanlan.zhihu.com/p/182553920`