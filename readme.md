# Codeforces-ContestCodeDownload

## Q & A

<img align="center" src="img/pixiv-86374864-small.jpg" />

Q：这个工具有什么用？

A：他能帮助你下载一场比赛中所有正式参赛选手的代码。

Q：官方不是支持了导出代码功能吗？

A：是的，但官方导出的代码是以`SubmissionID`命名的。在需要将下载的代码文件进行统计时这会很麻烦。特别是我在使用`SIM`或`Jplag`等工具进行代码查重时，我无法快速地看出重复率高的代码来自谁，需要一个个点进链接去确认。

Q：那官方导出代码也基本够用了吧？

A：不，还有一个问题。在我尝试管理我的Group里其他管理员创建的比赛时，我无法导出比赛代码（截图展示在最下方）。我不清楚这是不是一个特例，但好像只有比赛的创建者才能进行代码导出操作。因此，这一小工具也是为了让所有的管理员都能下载比赛代码进行分析，归档。

Q：那我该怎么使用呢？

A：你只需要将`release`中的`main.exe`文件放到一个合适的目录，再创建一个名叫`api.key`的文件在同一目录下。此后，你需要在`api.key`里面配置你在`Codeforces`的`API KEY`，`API SECRET`，`USERNAME`，`PASSWORD`。最后，你只需要双击运行`main.exe`，输入对应的比赛编号就可以了，比赛的代码文件会自动下载到同一目录。如果你还没有`API KEY`或`API SECRET`，你可以在[Settings - Codeforces](https://codeforces.com/settings/api)里创建一个。

Q：比赛的代码文件会以什么方式命名便于查看？

A：目前的命名方式是：`题号-题目名称-选手名称-所用语言(比赛ID#提交ID)`。需要注意的是，我只特殊处理了C/C++/Java/Python这四类使用的最多的语言，其他语言提交的代码会以`.txt`结尾，并在`所用语言`上标记为`Other`。（截图展示在最下方）

Q：好像没什么问题了...

A：其实还有点小问题。这次的核心代码用`Go`进行编写，因为初上手`Go`，而且时间较赶，其中的异常处理部分，日志部分，和数据库部分都还没有完成。因此`DEBUG`可能比较辛苦...

此外，未添加多线程，下载速度可能有限。

## 截图展示

#### Manager, can export submissions.

![image-20220705003842727](img/pic2.png)

---



#### Manager, but can't export submissions.

![image-20220705003836572](img/pic1.png)

---



#### Need to further confirm the relationship between the user and the submitted code.

![image-20220705010230911](img/jplag1.png)

---



#### The relationship between the user and the submitted code is quickly confirmed.

![image-20220705010601261](img/jplag2.png)

---



#### Presentation

![pre.gif](img/Pre.gif)


