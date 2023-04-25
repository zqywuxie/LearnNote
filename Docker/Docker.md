

# Docker

## 1.虚拟化和容器技术

---

###  **虚拟化技术**

> 虚拟化技术是一种将计算机物理资源进行抽象、转换为虚拟的计算机资源提供给程序使用的技术。

​	这里所指的计算机资源，就包括了 CPU 提供的**运算控制资源**，硬盘提供的**数据存储资源**，网卡提供的**网络传输资源**等。



####  **跨平台兼容**

---



​	从一个系统迁移到另一个系统，常会便随着程序的不兼容问题。那么虚拟化就是为了解决这个问题而诞生。在计算机技术发展的早期，各类计算平台、计算资源所提供的接口、调用方式十分杂乱，没有像今天这样相对统一的标准。**程序对计算机资源的调用主要依赖于操作系统所给出的接口。我们的程序通过操作系统提供的接口，向物理硬件发送指令。**

​	所以实现程序跨平台兼容方法很简单，**只要操作系统或者物理硬件所提供的接口调用方式一致即可**，程序便不需要兼容不同硬件平台的接口，只需要针对这一套接口开发即可。虚拟化技术正是通过其本身适配不同平台的硬件，而加以抽象成统一的接口，来实现程序跨平台运行这一目的的。



#### **虚拟化用于资源管理**

---



​	为应用程序设置一些虚假的资源数据，例如，我们只要告诉程序计算机只有 4GB 内存，那么不管真实的物理机是 8GB、16GB 还是 32GB，应用程序都会按照 4GB 这个虚假的值来处理它的逻辑。如此通过虚拟化技术管理计算机资源，可以对资源控制变得灵活并且**提高了计算机资源的使用率**。

​	这里提到了**提高计算机资源使用率**，可以使用虚拟化将原来程序用不到的资源，分享给另外程序，让资源不浪费。

​	例如，这里我们有一台运行 Nginx 的机器，由于 Nginx 运行对系统资源的消耗并不高，这就让系统几乎 95% 以上的资源处于闲置状态。这时候我们通过虚拟化技术，把其他的一些程序放到这台机器上来运行，它们就能够充分利用闲置的资源。这带来的好处就是我们不需要再为这些程序单独部署机器，从而节约不少的成本。

​	**问题**：我们本身可以在操作系统进行运行这些程序，为什么还要装到不同的虚拟环境？

​	**解答：**我们固然可以这么做，但要注意程序之间不会冲突。eg：端口用了同一个；不同程序依赖某个不同版本的工具库等。**虚拟化技术将资源进行了隔离，你用你的，我用我的。**那么就不存在冲突等问题了。



#### 虚拟化的分类

两大类：**硬件虚拟化，软件虚拟化**

- 硬件虚拟化

  物理硬件本身提供虚拟化；某个平台CPU能将其他平台的指令集转换为自身的指令集进行使用，给程序完全运行在那个平台的感觉。或者CPU自身模拟裂变，让操作系统或者软件认为存在多个CPU，进而同时运行多个程序或者操作系统。

- 软件虚拟化

  在软件虚拟化实现中，通过一层夹杂在应用程序和硬件平台上的虚拟化实现软件来进行指令的转换。也就是说，虽然应用程序向操作系统或者物理硬件发出的指令不是当前硬件平台所支持的指令，这个实现虚拟化的软件也会将之转换为当前硬件平台所能识别的。（wine)

  

当然，在实际场景中，虚拟化还能进行更加细化的分类，例如：

- **平台虚拟化**：在操作系统和硬件平台间搭建虚拟化设施，使得整个操作系统都运行在虚拟后的环境中。
- **应用程序虚拟化**：在操作系统和应用程序间实现虚拟化，只让应用程序运行在虚拟化环境中。
- **内存虚拟化**：将不相邻的内存区，甚至硬盘空间虚拟成统一连续的内存地址，即我们常说的虚拟内存。
- **桌面虚拟化**：让本地桌面程序利用远程计算机资源运行，达到控制远程计算机的目的



#### **虚拟机**

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/8/31/16590358c6f2217e~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

​	通常来说就是通过一个**虚拟机监视器 ( Virtual Machine Monitor )** 的设施来隔离操作系统与硬件或者应用程序和操作系统，以此达到虚拟化的目的。这个夹在其中的虚拟机监视器，常常被称为 **Hypervisor**。

​	这样使得程序或者操作系统可以无修改运行在另一平台上。但存在缺陷：**所有指令都需要经过`Hypervisor`转化，性能较为低下**。所以不完全遵循这种设计结构，会引入其他技术解决效率问题。

例如，在 VMware Workstation、Xen 中我们能够看到硬件辅助虚拟化的使用，通过让指令直达支持虚拟化的硬件，以此避开了效率低下的 Hypervisor。而如 JRE、HPHP 中，除了基于 Hypervisor 实现的**解释执行**机制外，还有**即时编译 ( Just In Time )** 运行机制，让程序代码在运行前编译成符合当前硬件平台的机器码，这种方式就已经不属于虚拟化的范畴了。





### 容器技术

---

​	所谓容器技术，指的是**操作系统自身支持一些接口**，能够让应用程序间可以互不干扰的独立运行，并且能够对其在运行中所使用的资源进行干预。也是属于操作系统虚拟化的范畴。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/1/1659296247facf28~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

​	简单来说就应用程序的运行被隔离在一个独立的环境，这个环境就像一个容器，包裹住应用程序，也正是容器名字的由来；

​	容器技术的**优势**：其在运行性能上要远超虚拟机等其他虚拟化实现。更甚一步说，运行在容器虚拟化中的应用程序，在运行效率上与真实运行在物理平台上的应用程序不相上下。

​	**原因：**因为**容器没有进行指令的转换**，由上可知虚拟化的效率低下主要是指令的转换，而容器技术却没有这一步。所以容器内部的应用程序必须支持在真是操作系统上运行，遵循硬件平台的指令规则。



#### **容器 VS 虚拟机**

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/1/16592899b28d4181~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

可以看到容器少了虚拟操作系统和虚拟机监视器两个层次，大幅度减少了资源消耗。

更准确的来说，**所有在容器中的应用程序其实完全运行在了宿主操作系统中**，与其他真实运行在其中的应用程序在指令运行层面是完全没有任何区别的。



## 2.Docker介绍

---

###  出现背景

​	一套系统的搭建一开始是基于本地的配置和依赖等，但将系统交给其他开发者或者测试，他们想要运行该系统就需要在本地搭建与你一样的配置。那么就会极大的降低工作效率和成本，开发者主要是进行实际开发，而不是纠缠在运行环境的问题。所以使用虚拟化技术进行优化。

![分布式应用服务体系](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/2/165997343db35f56~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)



#### **效率改变**

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/2/1659994bbc4225dd~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)



​	由上图就可得docker带来了巨大的效率提升。



### Docker技术实现

---

#### **三大技术：**

命名空间，控制组，联合文件系统。这些都是linux内核中的一些模块。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/2/16599a9d7a391ecf~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)



##### **命名空间**

​	目的是为了集合相同模块的类，区分不同模块间的同名类。类似分组，前端后端分一组，前端里面有个张三，后端也有个张三，但他们是两个互不干扰。

##### 控制组

​	资源控制组的作用就是控制计算机资源的。与以隔离进程、网络、文件系统等虚拟资源为目的 Namespace 不同，CGroups 主要做的是**硬件资源的隔离**。CGroups 除了资源的隔离，还有资源**分配**这个关键性的作用。通过 CGroups，我们可以指定任意一个隔离环境对任意资源的占用值或占用率，这对于很多分布式使用场景来说是非常有用的功能。

##### 联合文件系统

​	联合文件系统 ( Union File System ) 是一种能够同时挂载不同实际文件或文件夹到同一目录，形成一种**联合文件结构**的文件系统。与虚拟化本身无太大关系，而这里引入是为了解决文件系统占用过量，使得虚拟环境快速启停等问题。

​	在 Docker 中，提供了一种对 UnionFS 的改进实现，也就是 AUFS ( Advanced Union File System )。**解释：**将文件更新挂载到老的文件上，只改变更新的内容，不修改不更新的文件。eg：git提交，仓库不会全部重新更改，只会改变你提交的部分。



### Docker理念

---

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/7/165b2a9bd4a1a1b4~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

Docker 推崇一种轻量级容器的结构，即**一个应用一个容器**。

举个具体的例子，在常见的虚拟机实现中，我们要搭建一套 LAMP 结构的服务，我们通常会建立一个虚拟机，在虚拟机中安装上 Linux 系统，之后分别安装 Apache、MySQL 和 PHP。而在 Docker 里，**最佳的实践是分别基于 Apache、MySQL 和 PHP 的镜像建立三个容器，分别运行 Apache、MySQL 和 PHP ，而它们所在的虚拟操作系统也直接共享于宿主机的操作系统**。

如果我们将 Docker 的轻量级容器实现和虚拟机的一些参数进行对比，更容易得到结果。

| 属性           | Docker   | 虚拟机 |
| -------------- | -------- | ------ |
| 启动速度       | 秒级     | 分钟级 |
| 硬盘使用       | MB 级    | GB 级  |
| 性能           | 接近原生 | 较低   |
| 普通机器支撑量 | 数百个   | 几个   |



### Docker能做什么

---

#### 更快、更一致的交付你的应用程序

​	使用Docker后，能在本地容器中得到一套标准的应用或服务的运行环境，由此可以简化开发的生命周期 ( 减少在不同环境间进行适配、调整所造成的额外消耗 )。

#### 跨平台部署和动态伸缩

​	基于容器技术的 Docker 拥有很高的跨平台性，Docker 的容器能够很轻松的运行在开发者本地的电脑，数据中心的物理机或虚拟机，云服务商提供的云服务器，甚至是混合环境中。

**只要系统架构一样，是可以使用相同的镜像的**，比如x86的镜像只能x86的系统使用，arm的镜像只能arm系统使用。docker镜像对容器而言只是模拟了一个环境，跟宿主机没多大关系

​	同时，Docker 的**轻量性和高可移植性**能够很好的帮助我们完成应用的动态伸缩，我们可以通过一些手段近实时的对基于 Docker 运行的应用进行弹性伸缩，这能够大幅提高应用的健壮性。

#### 让同样的硬件提供更多的产出能力

​	Docker 的高效和轻量等特征，为替代基于 Hypervisor 的虚拟机提供了一个经济、高效、可行的方案。在 Docker 下，你能节约出更多的资源投入到业务中去，让应用程序产生更高的效益。同时，如此低的资源消耗也说明了 Docker 非常**适合在高密度的中小型部署场景中使用**。



### Docker核心组成

---

#### 四大组成对象

在 Docker 体系里，有四个对象 ( Object ) 是我们不得不进行介绍的，因为几乎所有 Docker 以及周边生态的功能，都是围绕着它们所展开的。它们分别是：**镜像 ( Image )**、**容器 ( Container )**、**网络 ( Network )**、**数据卷 ( Volume )**。



##### 镜像

镜像，可以理解为一个只读的文件包，其中包含了**虚拟环境运行最原始文件系统的内容**。

而Docker的镜像与虚拟机的镜像有区别。Docker 中的一个创新是利用了 **AUFS** 作为底层文件系统实现，通过这种方式，Docker 实现了一种**增量式的镜像结构**。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/7/165b29cad1a3dfae~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

每次对镜像内容的修改，Docker 都会将这些**修改铸造成一个镜像层**，而一个镜像其实就是由其下层所有的镜像层所组成的。当然，每一个镜像层单独拿出来，与它之下的镜像层都可以组成一个镜像。

另外，由于这种结构，Docker 的镜像**实质上是无法被修改**的，因为所有对镜像的修改只会产生新的镜像，而不是更新原有的镜像。



##### 容器

容器技术中，容器就是用来隔离虚拟环境的基础设施，而在 Docker 里，它也被引申为**隔离出来的虚拟环境**。

如果把**镜像理解为编程中的类**，那么**容器就可以理解为类的实例**。镜像内存放的是不可变化的东西，当以它们为基础的容器启动后，容器内也就成为了一个“活”的空间。

Docker 的容器应该有三项内容组成：

- 一个 Docker 镜像
- 一个程序运行环境
- 一个指令集合

##### 网络

前面说了容器是相互隔离的，但要与外界或者其他程序进行交互，这里指的交互大多情况指的数据信息的交互。网络交互就是目前最常用的一个程序间的数据交互方式。

由于计算机网络领域拥有相对统一且独立的协议等约定，其跨平台性非常优秀，所有的应用都可以通过网络在不同的硬件平台或操作系统平台上进行数据交互。

Docker实现网络功能后，可以对每个容器进行网络配置，在容器之间建立虚拟网络，与其他网络环境形成隔离。如下图

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/5/165a810ad2c81714~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

另外，利用一些技术，Docker 能够在**容器中营造独立的域名解析环境**，这使得我们可以在不修改代码和配置的前提下直接迁移容器，**Docker 会为我们完成新环境的网络适配**。对于这个功能，我们甚至能够在不同的物理服务器间实现，让处在两台物理机上的两个 Docker 所提供的容器，加入到同一个虚拟网络中，形成完全屏蔽硬件的效果。



##### 数据卷

卷就是目录或文件，存在于一个或多个容器中，由Docker挂载到容器，但卷不属于联合文件系统（Union FileSystem），因此能够绕过联合文件系统提供一些用于持续存储或共享数据的特性:。

**卷的设计目的就是数据的持久化，完全独立于容器的生存周期，因此Docker不会在容器删除时删除其挂载的数据卷。**



#### Docker Engine

---



Docker Engine是Docker中最核心的软件，在 Docker Engine 中，实现了 Docker 技术中最核心的部分，也就是容器引擎这一部分



##### docker daemon 和 docker CLI

Docker Engine也是由多个独立软件所组成的软件包。其中最核心的激素**docker daemon（无交互后台程序) 和 docker CLI(命令行界面）**

所有我们通常认为的 Docker 所能提供的容器管理、应用编排、镜像分发等功能，都集中在了 docker daemon 中,前面的四大组成对象也都是现在其中。

在操作系统里，docker daemon 通常以服务的形式运行以便静默的提供这些功能，所以我们也通常称之为 **Docker 服务。**

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/5/165a8349ffdb33e0~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

docker daemon在后台管理容器相关资源的同时，也向外暴露了一套RESTful API，用户通过这些接口进行操控docker daemon所管理的相关资源。docker也就提供了docker CLI来帮助我们对接口的请求，也就不需要自己编写HTTP请求。



## 3.搭建环境

---

这里以centos7演示，因为docker对主流Linux系统有一些版本要求；

### 安装（Linux）

```sh
$ sudo yum install yum-utils device-mapper-persistent-data lvm2
$
$ sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
$ sudo yum install docker-ce
$
$ sudo systemctl enable docker #自启动
$ sudo systemctl start docker # 启动
```

**安装好查看版本**

```sh
$ sudo docker version
Client: Docker Engine - Community
 Version:           23.0.1
 API version:       1.42
 Go version:        go1.19.5  #docker基于go语言
 Git commit:        a5ee5b1
 Built:             Thu Feb  9 19:51:00 2023
 OS/Arch:           linux/amd64
 Context:           default
permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get "http://%2Fvar%2Frun%2Fdocker.sock/v1.24/version": dial unix /var/run/docker.sock: connect: permission denied
```

**查看更多信息**

```sh
$ sudo docker info
Client:
 Context:    default
 Debug Mode: false
 Plugins:
  buildx: Docker Buildx (Docker Inc.)
    Version:  v0.10.2
    Path:     /usr/libexec/docker/cli-plugins/docker-buildx
  compose: Docker Compose (Docker Inc.)
    Version:  v2.16.0
    Path:     /usr/libexec/docker/cli-plugins/docker-compose
  scan: Docker Scan (Docker Inc.)
    Version:  v0.23.0
    Path:     /usr/libexec/docker/cli-plugins/docker-scan

Server:
 Containers: 0
  Running: 0
  Paused: 0
  Stopped: 0
 Images: 0
 Server Version: 23.0.1
....
 Live Restore Enabled: false
```

可以看到正在运行的 Docker Engine 实例中运行的容器数量，存储的引擎等等信息。



**配置国内镜像源**

Docker 中也有一个由官方提供的中央镜像仓库。不过国外站点你懂的，这里使用官方提供的国内镜像源

> registry.docker-cn.com

修改文件`/etc/docker/daemon.json`,如果不存在你可以创建它

```sh
{
    "registry-mirrors": [
        "https://registry.docker-cn.com"
    ]
}
```

记得重启docker daemon

```sh
$ sudo systemctl restart docker
```

`docker info`进行查看当前镜像源

```sh
.....
Registry: https://index.docker.io/v1/
 Experimental: false
 Insecure Registries:
  127.0.0.0/8
 Registry Mirrors:
  https://registry.docker-cn.com/
 Live Restore Enabled: false

```

读者切记linux下使用管理员身份查看信息，否则看不到



### 安装（windows/Mac）

[Docker for Windows](https://link.juejin.cn/?target=https%3A%2F%2Fstore.docker.com%2Feditions%2Fcommunity%2Fdocker-ce-desktop-windows)

[Docker for Max](https://link.juejin.cn/?target=https%3A%2F%2Fstore.docker.com%2Feditions%2Fcommunity%2Fdocker-ce-desktop-mac)

windows需要先安装wsl，然后可以启动Docker Desktop。

安装后启动即可在任务栏看到大鲸鱼。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/10/165c1d1fb7030b63~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

一直闪动说明在部署docker daemon相关的配置和启动。当不山东说明就可以使用了。Docker Desktop 为我们在 Windows 和 macOS 中使用 Docker 提供了与 Linux 中几乎一致的方法，所以在命令行同上linux命令执行即可。

```sh
λ docker version
Client:
## ......
 OS/Arch:  windows/amd64
## ......
```



### Docker Desktop 的实现原理

---

前面讲到docker的容器实现是基于Linux内核的三大技术等功能的。那么为什么能使用docker呢。

因为Windows和macos本身也具有虚拟化的功能，Docker for Windows 和 Docker for Mac 这里利用了这两个操作系统提供的功能来搭建一个虚拟 Linux 系统，并在其之上安装和运行 docker daemon。Windows里面的便是**wsl（适用于Linux的Windows子系统）**

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/12/165cb3b94b24b951~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

#### 主机文件挂载

Docker 容器中能够通过数据卷的方式挂载宿主操作系统中的文件或目录，宿主操作系统在 Windows 和 macOS 环境下的 Docker Desktop 中，指的是虚拟的 Linux 系统。但我们期望的是挂载Windows和macOS里面文件。

而实现这个效果，我们可以将目录挂载在虚拟的Linux系统上，然后再用docker挂载到容器中，整个过程就被集合在Docker Desktop中，不需要人工操作，实现了自动化。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/11/165c8400bf8f809e~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)



相关使用和配置需要就进行百度，这里不多阐述；

## 4.镜像和容器

### Dcoker镜像

---



可以将 Docker 镜像理解为包含应用程序以及其相关依赖的一个基础文件系统，在 Docker 容器启动的过程中，它**以只读的方式被用于创建容器的运行环境**。



#### 深入镜像实现

与其他虚拟机的镜像管理不同，Docker 将镜像管理纳入到了自身设计之中，也就是说，所有的 Docker 镜像都是按照 Docker 所设定的逻辑打包的，也是受到 Docker Engine 所控制的。

例如我们常见的虚拟机镜像，通常是由热心的提供者以他们自己熟悉的方式打包成镜像文件，被我们从网上下载或是其他方式获得后，恢复到虚拟机中的文件系统里的。而 Docker 的镜像我们**必须通过 Docker 来打包，也必须通过 Docker 下载或导入后使用**，不能单独直接恢复成容器中的文件系统。

虽然这么做失去了很多灵活性，但固定的格式意味着我们可以很轻松的在不同的服务器间传递 Docker 镜像，配合 Docker 自身对镜像的管理功能，让我们在不同的机器中传递和共享 Docker 变得非常方便。这也是 Docker 能够提升我们工作效率的一处体现。

对于每一个记录文件系统修改的镜像层来说，Docker 都会根据它们的信息生成了一个 Hash 码，**这是一个 64 长度的字符串**，足以保证全球唯一性。这种编码的形式在 Docker 很多地方都有体现，之后我们会经常见到。

由于镜像层都有唯一的编码，我们就能够区分不同的镜像层并能保证它们的内容与编码是一致的，这带来了另一项好处，就是**允许我们在镜像之间共享镜像层**。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/13/165d0692fe7a478b~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

比如两个镜像都是基于一个相同的镜像修改而得，那么在实际使用中，就可以进行共享该镜像内部的镜像层。

好处可以使得镜像共享一些存储空间，达到1 + 1 < 2的效果。

#### 查看镜像

`docker images`查看镜像信息

```sh
$ docker images
REPOSITORY    TAG       IMAGE ID       CREATED         SIZE
redis         latest    f9c173b0f012   12 days ago     117MB
hello-world   latest    feb5d9fea6a5   17 months ago   13.3kB
```

我们发现在结果中镜像 ID 的长度只有 12 个字符，这和我们之前说的 64 个字符貌似不一致。其实为了避免屏幕的空间都被这些看似“乱码”的镜像 ID 所挤占，**所以 Docker 只显示了镜像 ID 的前 12 个字符**，大部分情况下，它们已经能够让我们在单一主机中识别出不同的镜像了。

#### 镜像命名

虽然镜像ID可以识别出镜像，但是这样的命名太长不合理。所以需要对镜像进行一个命令

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/12/165cc15252cc5e51~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

准确的来说，镜像的命名我们可以分成三个部分：**username**、**repository** 和 **tag**。

- **username**： 主要用于识别上传镜像的不同用户，与 GitHub 中的用户空间类似。
- **repository**：主要用于**识别进行的内容**，形成对镜像的表意描述。
- **tag**：主要用户表示**镜像的版本**，方便区分进行内容的不同细节

但上面展示`docker images`，有些并没有username这个部分，表示镜像是由 Docker 官方所维护和提供的，所以就不单独标记用户了。

并且可以看`REPOSITORY` 常用软件名表示。但与软件命名是分开的，之所以采用软件名，在于docker对容器的轻量级设计。通常会只在一个容器中运行一个应用程序，那么自然容器的镜像也会仅包含程序以及与它运行有关的一些依赖包，所以我们使用程序的名字直接套用在镜像之上，既祛除了镜像取名的麻烦，又能直接表达镜像中的内容。

`tag`命名也通常参考镜像所关联的应用程序，比如版本号等环境，构建方式信息。`php:7.2-cli` 和 `php:7.2-fpm` 就包含了镜像的构建方式和版本。如果没有给出具体的tag，那么就采用`latest`如上显示，保持软件最新版的使用。



### 容器的生命周期

---

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/17/165e53743e730432~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

着重看看容器的几个核心状态，也就是图中色块表示的：**Created**、**Running**、**Paused**、**Stopped**、**Deleted**。

#### 主进程

在 Docker 的设计中，容器的生命周期其实与容器中 PID 为 1 这个进程有着密切的关系。容器启动理解为该进程启动，停止也意味着进程停止。

虽然在 Docker 中我们也能够实现在同一个容器中运行多个不同类型的程序，但这么做的话，Docker 就无法跟踪不同应用的生命周期，有可能造成应用的非正常关闭，进而影响系统、数据的稳定性。



### 写时复制机制

---



在编程里，写时复制常常用于对象或数组的拷贝中，当我们拷贝对象或数组时，**复制的过程并不是马上发生在内存中**，而只是先让两个变量同时指向同一个内存空间，并进行一些标记，当我们要对对象或数组进行修改时，才真正进行内存的拷贝。

当 Docker 第一次启动一个容器时，**初始的读写层是空的**，当文件系统发生变化时，这些变化都会应用到这一层之上。比如，如果想修改一个文件，这个文件首先会从该读写层下面的只读层复制到该读写层。由此，该文件的只读版本依然存在于只读层，只是被读写层的该文件副本所隐藏。该机制则被称之为**写时复制（Copy on write）**。



## 5.获得镜像

---

前面提到一个`镜像仓库`概念，可以理解为GitHub等平台。主要好处并不止存储镜像，而是镜像的分发。开发环境下进行推送，然后在测试或者生产环境进行拉取，几个命令甚至自动化就可完成。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/18/165eacb6b1b2c1ac~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

### 拉取镜像

`docker pull`

```sh
$ sudo docker pull ubuntu
Using default tag: latest
latest: Pulling from library/ubuntu
124c757242f8: Downloading [===============================================>   ]  30.19MB/31.76MB
9d866f8bde2a: Download complete 
fa3f2f277e67: Download complete 
398d32b153e8: Download complete 
afde35469481: Download complete 
```

可以看到下载分几行，每一行就代表一个镜像层。拉取每一个镜像层然后组合为镜像。如果本地有相同镜像层，就会略过这个下载，如前面所讲共享镜像层。

```sh
$ docker images
REPOSITORY    TAG       IMAGE ID       CREATED         SIZE
redis         latest    f9c173b0f012   12 days ago     117MB
ubuntu        latest    74f2314a03de   13 days ago     77.8MB
hello-world   latest    feb5d9fea6a5   17 months ago   13.3kB
```

### Docker hub

docker官方建立中央镜像仓库，类似maven仓库。然后搜索自己想要的镜像

除了网站搜索，还可以使用`docker search`

```sh
$ sudo docker search ubuntu
NAME                                                   DESCRIPTION                                     STARS               OFFICIAL            AUTOMATED
ubuntu                                                 Ubuntu is a Debian-based Linux operating sys…   8397                [OK]                
dorowu/ubuntu-desktop-lxde-vnc                         Ubuntu with openssh-server and NoVNC            220                                     [OK]
rastasheep/ubuntu-sshd                                 Dockerized SSH service, built on top of offi…   171                                     [OK]
consol/ubuntu-xfce-vnc                                 Ubuntu container with "headless" VNC session…   129                                     [OK]
ansible/ubuntu14.04-ansible                            Ubuntu 14.04 LTS with ansible                   95                                      [OK]
ubuntu-upstart                                         Upstart is an event-based replacement for th…   89                  [OK]                
neurodebian                                            NeuroDebian provides neuroscience research s…   54                  [OK]                
## ......
```



### 管理镜像

---

获得镜像更详细的信息，`docker inspect`

```sh
$ sudo docker inspect redis:3.2
[
    {
        "Id": "sha256:2fef532eadb328740479f93b4a1b7595d412b9105ca8face42d3245485c39ddc",
        "RepoTags": [
            "redis:3.2"
        ],
        "RepoDigests": [
            "redis@sha256:745bdd82bad441a666ee4c23adb7a4c8fac4b564a1c7ac4454aa81e91057d977"
        ],
## ......
    }
]
```

还可以传入容器ID

```sh
$ sudo docker inspect redis:4.0
$ sudo docker inspect 2fef532e
```



这里可以看的容器ID我们只传入了8位，其实传入1位也可以，前提是可以找到唯一的镜像

```sh
C:\Users\wuxie>docker inspect 6
[
    {
        "Id": "610dff2378ded2940babe360c78c57f2f96f8d260164f1719d2135005366b323",
        "Created": "2023-03-13T12:45:25.333206594Z",
        "Path": "docker-entrypoint.sh",
        "Args": [
            "--requirepass",
```

如果找不到唯一的镜像，那么就报错



#### 删除镜像

`docker rmi id/name`  

```sh
$ docker rmi ubuntu
Untagged: ubuntu:latest
Untagged: ubuntu@sha256:2adf22367284330af9f832ffefb717c78239f6251d9d0f58de50b86229ed1427
Deleted: sha256:74f2314a03de34a0a2d552b805411fc9553a02ea71c1291b815b2f645f565683
Deleted: sha256:202fe64c3ce39b94d8beda7d7506ccdfcf7a59f02f17c915078e4c62b5c2ed11
```

删除镜像也是删除镜像内的镜像层，也并不会删除共享镜像层。如果某个镜像不存在多个标签，当且仅当只有一个标签时，执行删除命令时，会彻底删除镜像。

也可以删除多个镜像

```sh
$ sudo docker rmi redis:3.2 redis:4.0
Untagged: redis:3.2
Untagged: redis@sha256:745bdd82bad441a666ee4c23adb7a4c8fac4b564a1c7ac4454aa81e91057d977
Deleted: sha256:2fef532eadb328740479f93b4a1b7595d412b9105ca8face42d3245485c39ddc
## ......
Untagged: redis:4.0
Untagged: redis@sha256:b77926b30ca2f126431e4c2055efcf2891ebd4b4c4a86a53cf85ec3d4c98a4c9
Deleted: sha256:e1a73233e3beffea70442fc2cfae2c2bab0f657c3eebb3bdec1e84b6cc778b75
## ......
```



## 6.运行和管理容器

### 容器的创建和启动

前面讲了容器的生命周期

- **Created**：容器已经被创建，容器所需的相关资源已经准备就绪，但容器中的程序还未处于运行状态。
- **Running**：容器正在运行，也就是容器中的应用正在运行。
- **Paused**：容器已暂停，表示容器中的所有程序都处于暂停 ( 不是停止 ) 状态。
- **Stopped**：容器处于停止状态，占用的资源和沙盒环境都依然存在，只是容器中的应用程序均已停止。
- **Deleted**：容器已删除，相关占用的资源及存储在 Docker 中的管理信息也都已释放和移除。

#### 创建容器

选择好镜像 `docker create image`

```sh
$ sudo docker create nginx:1.12
34f277e22be252b51d204acbb32ce21181df86520de0c337a835de6932ca06c3
```

可以给容器命名

```sh
$ sudo docker create --name nginx nginx:1.12
```



#### 启动容器

通过 `docker create` 创建的容器，是处于 Created 状态的，其内部的应用程序还没有启动，所以我们需要通过 `docker start` 命令来启动它。

```sh
$ sudo docker start nginx
```

通过 `docker run` 这个命令将 `docker create` 和 `docker start` 这两步操作合成为一步，进一步提高工作效率。

```sh
$ sudo docker run --name nginx -d nginx:1.12
89f2b769498a50f5c35a314ab82300ce9945cbb69da9cda4b022646125db8ca7
```

`docker run`会采用去前台运行，导致控制台阻塞了，所以可以使用`-d`或`-detach` 使得程序在后台运行。

#### 容器相关命令

`docker rename 原名 新名` 修改容器名字

`docker ps` -l 最近创建的容器，-n x 显示最近x个创建的容器**，-q静默模式，只显示容器ID**

```sh
$ docker run -it ubuntu bash
# ctrl+p+q 退出容器，但不停止
# exit ctrl+c/d 退出并停止
```

也可以create,start 然后直接进入容器



docker kill ID/name 强制停止容器

docker rm -f xx 强制删除正在运行的容器

```sh
docker run --name redis-6667 -p 6667:6379 -v /E/DockerConfig/redis/conf/redis.conf:/etc/redis/redis_6379.conf -v /E/DockerConfig/redis/data:/data/ -d redis:latest redis-server /etc/redis/redis_6379.conf --appendonly yes
```



### 管理容器

---

`docker ps` 列举出docker中运行的容器。`-a/-all`列出所有状态的容器

`ps`列出容器的过程是查看正在运行的进程，所以才使用同Linux中类似ps的指令，而不是ls。

```sh
$ docker ps -a
CONTAINER ID   IMAGE          COMMAND                  CREATED          STATUS                      PORTS     NAMES
0fe8511f1039   hello-world    "/hello"                 38 minutes ago   Exited (0) 38 minutes ago             infallible_bartik
610dff2378de   redis:latest   "docker-entrypoint.s…"   19 hours ago     Exited (0) 19 hours ago               redis-80wo
```

- CONTAINER ID 容器id
- image 容器基于的镜像
- CREATED 容器的创建时间
- names 容器的名称
- COMMAND 程序的启动命令
- STATUS 容器所处状态
  - **Created** 此时容器已创建，但还没有被启动过。
  - **Up [ Time ]** 这时候容器处于正在运行状态，而这里的 Time 表示容器从开始运行到查看时的时间。
  - **Exited ([ Code ]) [ Time ]** 容器已经结束运行，这里的 Code 表示容器结束运行时，主程序返回的程序退出码，而 Time 则表示容器结束到查看时的时间。



#### 停止和删除容器

`docker stop` 停止容器

```sh
$ sudo docker stop nginx
```

`docker rm`删除容器

```sh
$ sudo docker rm nginx
```

正在运行的容器是不能被删除的，可以增加`-f/-force`进行强制停止并删除，但不推荐这种做法；

#### 随手删除容器

与其他虚拟机不同，Docker 的轻量级容器设计，讲究随用随开，随关随删。也就是说，当我们短时间内不需要使用容器时，最佳的做法是删除它而不是仅仅停止它。

有的读者会问，容器一旦删除，其内部的文件系统变动也就消失了，这样做岂不是非常麻烦。要解决这个疑惑，其根本是解决为什么我们会对容器中的文件系统做更改。我这里总结了两个对虚拟环境做更改的原因，以及在 Docker 中如何优雅的解决它们。

- 在使用虚拟机或其他虚拟化所搭建的虚拟环境时，我们倾向于使用一个干净的系统镜像并搭建程序的运行环境，由于将这类虚拟环境制作成镜像的成本较高，耗时也非常久，所以我们对于一些细小的改动倾向于修改后保持虚拟环境不被清除即可。而在 Docker 中，打包镜像的成本是非常低的，其速度也快得惊人，**所以如果我们要为程序准备一些环境或者配置，完全可以直接将它们打包至新的镜像中，下次直接使用这个新的镜像创建容器即可。**
- 容器中应用程序所产生的一些文件数据，是非常重要的，如果这些数据随着容器的删除而丢失，其损失是非常巨大的。对于这类由应用程序所产生的数据，并且需要保证它们不会随着容器的删除而消失的，**我们可以使用 Docker 中的数据卷来单独存放**。由于数据卷是独立于容器存在的，所以其能保证数据不会随着容器的删除而丢失。关于数据卷的具体使用，在之后的小节会专门讲解。

解决了这两个问题，大家心中的疑虑是不是就小了很多。而事实上，容器的随用随删既能保证在我们不需要它们的时候它们不会枉占很多资源，也保证了每次我们建立和启动容器时，它们都是“热乎”的崭新版本。大家都知道，系统卡就重装，而借助 Docker 秒级的容器启停特性，我们就是可以这么任性的“重装”。



### 进入容器

---

`docker exec` 命令能帮助我们在正在运行的容器中运行指定命令。在Linux系统中我们通过控制台软件进行操控Linux，熟悉的是`shell/bash`，分别有sh和bash两个程序启动。

所以容器启动这两个程序就可以对容器进行操控。

bash比sh功能丰富，所以优选bash

```sh
$ sudo docker exec -it nginx bash
root@83821ea220ed:/#
```

两个选项不可或缺,其中 `-i` ( `--interactive` ) 表示保持我们的输入流，只有使用它才能保证控制台程序能够正确识别我们的命令。而 `-t` ( `--tty` ) 表示启用一个伪终端，形成我们与 bash 的交互，如果没有它，我们无法看到 bash 内部的执行结果。

#### 衔接到容器

`docker attach `用于将当前的输入输出流连接到指定的容器上。

简单来说就是将容器主程序转到了前台运行(`docker run -d`有相反意思)

由于我们的输入输出流衔接到了容器的主程序上，我们的输入输出操作也就直接针对了这个程序，而我们发送的 Linux 信号也会转移到这个程序上。例如我们可以通过 Ctrl + C 来向程序发送停止信号，让程序停止 ( 从而容器也会随之停止 )。

在实际开发中，由于 `docker attach` 限制较多，功能也不够强大，所以并没有太多用武之地，这里我们就一笔带过，不做详细的解读了。

> attach是可以带上--sig-proxy=false来确保CTRL-D或CTRL-C不会关闭容器。
> sudo docker attach --sig-proxy=false nginx







## 7.🥴管理和存储数据

### 数据管理实现方式

---

docker文件系统的缺点:

- **沙盒文件系统是跟随容器生命周期所创建和移除的**，数据无法直接被持久化存储。
- 由于容器隔离，我们很难从容器外部获得或操作容器内部文件中的数据。

因为docker容器文件系统基于UnionFS ，可以支持挂载不同类型的文件系统到统一的目录结构，所以只需要将宿主操作系统中的文件或目录挂载到容器中，就可以让容器内外共享这个文件。同时，UnionFS 带来的**读写性能损失是可以忽略不计**的；



#### 挂载方式

Docker 提供了三种适用于不同场景的文件系统挂载方式：**Bind Mount**、**Volume** 和 **Tmpfs Mount**。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/25/1660eff4b182c891~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

- **Bind Mount** 能够直接将宿主操作系统中的目录和文件挂载到容器内的文件系统中，通过指定容器外的路径和容器内的路径，就可以形成挂载映射关系，在容器内外对文件的读写，都是相互可见的。
- **Volume** 也是从宿主操作系统中挂载目录到容器内，只不过**这个挂载的目录由 Docker 进行管理**，我们只需要指定容器内的目录，不需要关心具体挂载到了宿主操作系统中的哪里。
- **Tmpfs Mount** 支持**挂载系统内存中的一部分到容器的文件系统里**，不过由于内存和容器的特征，它的**存储并不是持久**的，其中的内容会随着容器的停止而消失。



### 挂载文件到容器

---

我们可以在容器创建的时候通过传递 `-v` 或 `--volume` 选项来指定内外挂载的对应目录或文件

```sh
$ sudo docker run -d --name nginx -v /webapp/html:/usr/share/nginx/html nginx:1.12
```

```sh
-v <host-path>:<container-path>` 或 `--volume <host-path>:<container-path>
```

我们能够指定目录进行挂载，也能够指定具体的文件来挂载,根据实际情况而定。

挂载后就可以在容器中查看宿主操作系统中的文件了。

`docker inspect` 查看相关信息

```sh
$ sudo docker inspect nginx
[
    {
## ......
        "Mounts": [
            {
                "Type": "bind",
                "Source": "/webapp/html",
                "Destination": "/usr/share/nginx/html",
                "Mode": "",
               #支持读写
                "RW": true,
                "Propagation": "rprivate"
            }
        ],
## ......
    }
]
```



```sh
$ sudo docker run -d --name nginx -v /webapp/html:/usr/share/nginx/html:ro nginx:1.12
```

`:ro`就表示只读

注意由于权限问题可以挂载宿主操作系统任何目录和文件，对安全性造成了影响，所以使用`Bind Mount`需要注意选择，如下几种使用场景适合

- 当我们需要从**宿主操作系统共享配置**的时候。对于一些配置项，我们可以直接从容器外部挂载到容器中，这利于保证容器中的配置为我们所确认的值，也方便我们对配置进行监控。例如，遇到容器中时区不正确的时候，我们可以直接将操作系统的时区配置，也就是 /etc/timezone 这个文件挂载并覆盖容器中的时区配置。
- 当我们需要借助 Docker 进行开发的时候。虽然在 Docker 中，推崇直接将代码和配置打包进镜像，以便快速部署和快速重建。但这在开发过程中显然非常不方便，因为每次构建镜像需要耗费一定的时间，这些时间积少成多，就是对开发工作效率的严重浪费了。如果我们直接把**代码挂载进入容器，那么我们每次对代码的修改都可以直接在容器外部进行**。

关于绑定挂载的另一个重要信息是，它们可以访问敏感文件。根据Docker文档，你可以通过在容器中运行的进程改变主机文件系统。这包括创建、修改和删除系统文件和目录，这可能有相当严重的安全影响。它甚至可能影响到非Docker进程。





#### 临时挂载

`Tmpfs Mount` 主要利用内存来存储数据。由于内存不是持久性存储设备，所以其带给 Tmpfs Mount 的特征就是临时性挂载。

内存位置不需要我们指定，只需要挂载到容器内的目录即可,使用选项--tmpfs`。

```sh
$ sudo docker run -d --name webapp --tmpfs /webapp/cache webapp:latest
```

```sh
$ sudo docker inspect webapp
[
    {
## ......
         "Tmpfs": {
            "/webapp/cache": ""
        },
## ......
    }
]
```

挂载临时文件首先要注意它不是持久存储这一特性，在此基础上，它有几种常见的适应场景。

- 应用中使用到，但不需要进行持久保存的敏感数据，可以借助内存的非持久性和程序隔离性进行一定的安全保障。
- 读写速度要求较高，数据变化量大，但不需要持久保存的数据，可以借助内存的高读写速度减少操作的时间。



### 数据卷

---

数据卷本质是宿主操作系统上的一个目录，放在Docker内部，接受docker的管理。挂载时不需要知道存在宿主系统的何处。

```sh
$ sudo docker run -d --name webapp -v /webapp/storage webapp:latest
```

```sh
$ sudo docker inspect webapp
[
    {
## ......
        "Mounts": [
            {
                "Type": "volume",
                "Name": "2bbd2719b81fbe030e6f446243386d763ef25879ec82bb60c9be7ef7f3a25336",
                "Source": "/var/lib/docker/volumes/2bbd2719b81fbe030e6f446243386d763ef25879ec82bb60c9be7ef7f3a25336/_data",
                "Destination": "/webapp/storage",
                "Driver": "local",
                "Mode": "",
                "RW": true,
                "Propagation": ""
            }
        ],
## ......
    }
]
```

信息注意与绑定挂载区分。

- `type`显然不同
- `Name` 默认采用数据卷的ID命名。也可以自定义

```sh
# -v <name>:<container-path>

$ sudo docker run -d --name webapp -v appdata:/webapp/storage webapp:latest
```

- `Source` 分配用于挂载的宿主目录，默认`var/lib/docker`，不需要关注，docker管理。

> 注意绑定时`-v`选项，与前面绑定挂载区分。绑定挂载必须使用**绝对路径**，避免与数据卷挂载中命名这种形式冲突。
>



虽然与绑定挂载的原理差别不大，但数据卷在许多实际场景下你会发现它很有用。

- 当希望将数据在**多个容器间共享**时，利用数据卷可以在保证数据持久性和完整性的前提下，完成更多自动化操作。
- 当我们希望对容器中挂载的内容进行管理时，可以直接利用数据卷自身的管理方法实现（比挂载绑定安全)。
- 当使用远程服务器或云服务作为存储介质的时候，**数据卷能够隐藏更多的细节**，让整个过程变得更加简单。



#### 共用数据卷

数据卷用来实现容器间的目录共享，也可以使用绑定挂载实现，但是数据卷更为简单。

由于**数据卷的命名在 Docker 中是唯一的**，所以我们很容易通过数据卷的名称确定数据卷，这就让我们很方便的让多个容器挂载同一个数据卷了。



```shell
$ sudo docker run -d --name webapp -v html:/webapp/html webapp:latest
$ sudo docker run -d --name nginx -v html:/usr/share/nginx/html:ro nginx:1.12
```

> `privileged` 使用该参数，[container](https://so.csdn.net/so/search?q=container&spm=1001.2101.3001.7020)内的root拥有真正的root权限。
>
> Docker挂载主机目录访问如果出现cannot open directory .: Permission denied
>
> 解决办法：在挂载目录后多加一个--privileged=true参数即可

数据卷不存在，Docker就会自动创建和分配。

可以下面命令操作数据卷。

`docker volume create` 不依赖容器独立创建数据卷

```shell
$ sudo docker volume create appdata
```

通过 `docker volume ls` 可以列出当前已创建的数据卷。

```shell
$ sudo docker volume ls
DRIVER              VOLUME NAME
local               html
local               appdata
```



#### 继承容器卷

docker run -it --privileged=true --volumes-from 父类容器名 --name 子类容器名 镜像名

```sh
[root@docker ~]# docker run -it --privileged=true --volumes-from my-ubt --name my-ubt-son ubuntu
```

my-ubt-son是继承my-ubt的映射规则，因此即使my-ubt停止，宿主机的数据仍能同步到my-ubt-son 所以能看到之前共享数据卷内创建的文件

#### 删除数据卷

该用专门的命令进行删除，而不是直接去目录进行删除。

`docker volume`

```shell
$ sudo docker volume rm appdata
```

删除之前，确保数据卷没有被任何容器使用，否则Docker不允许删除。



对于我们没有直接命名的数据卷，因为要反复核对数据卷 ID，这样的方式并不算特别友好。这种没有命名的数据卷，通常我们可以看成它们与对应的容器产生了绑定，因为其他容器很难使用到它们。而这种绑定关系的产生，也让我们可以在容器删除时将它们一并删除。

在 `docker rm` 删除容器的命令中，我们可以通过增加 `-v` 选项来删除容器关联的数据卷。

```shell
$ sudo docker rm -v webapp
```



如果我们没有随容器删除这些数据卷，Docker 在创建新的容器时也不会启用它们，**即使它们与新创建容器所定义的数据卷有完全一致的特征**。

Docker 向我们提供了 `docker volume prune` 这个命令，它可以删除那些没有被容器引用的数据卷。

```shell
$ sudo docker volume prune -f
Deleted Volumes:
af6459286b5ce42bb5f205d0d323ac11ce8b8d9df4c65909ddc2feea7c3d1d53
0783665df434533f6b53afe3d9decfa791929570913c7aff10f302c17ed1a389
65b822e27d0be93d149304afb1515f8111344da9ea18adc3b3a34bddd2b243c7
## ......
```



### 数据卷容器

---

一个没有具体指定的应用，甚至不需要运行的容器，作用定义一个或多个数据卷并持有它们的引用。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/26/166135778cfd74c2~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)



创建数据卷容器的方式很简单，由于不需要容器本身运行，因而我们找个简单的系统镜像都可以完成创建。

```shell
$ sudo docker create --name appdata -v /webapp/storage ubuntu
```

在使用数据卷容器时，我们**不建议再定义数据卷的名称**，因为我们可以**通过对数据卷容器的引用来完成数据卷的引用**。而不设置数据卷的名称，也避免了在同一 Docker 中数据卷重名的尴尬。



数据卷容器可以作为容器文件系统桥梁，创建新容器时进行使用`--volumes-from`

```sh
$ sudo docker run -d --name webapp --volumes-from appdata webapp:latest
```

引用数据卷容器时，**不需要再定义数据卷挂载到容器中的位置**，Docker 会以数据卷容器中的挂载定义将数据卷挂载到引用的容器中。 隐藏了数据卷的配置和定制，而这些隐藏，意味着我们能够更轻松的实现容器的迁移。



#### 备份和迁移数据卷

由于数据卷本身就是宿主操作系统中的一个目录，我们只需要在 Docker 资源目录里找到它就可以很轻松的打包、迁移、恢复了。但尽量采用docker中的方法

**数据备份、迁移、恢复的过程可以理解为对数据进行打包，移动到其他位置**，在需要的地方解压的过程。在数据打包之前，我们先建立一个用来存放打包文件的目录，这里我们使用 `/backup` 作为例子。

要备份数据，我们**先建立一个临时的容器**，将用于备份的目录和要备份的数据卷都挂载到这个容器上。

```sh
$ sudo docker run --rm --volumes-from appdata -v /backup:/backup ubuntu tar cvf /backup/backup.tar /webapp/storage
```

- `--rm` 停止后自动删除
- `tar cvf /backup/backup.tar /webapp/storage` tar cvf 压缩，将/webapp/storage内的数据压缩至/backup 取名为backup.tar,

解压同上，压缩命令换成解压

```sh
$ docker run --rm --volumes-from appdata -v /backup:/backup ubuntu tar xvf /backup/backup.tar -C /webapp/storage --strip
```

- `-C` 指定目录
- `--strip`

### 另一个挂载选项

---

`-v`选项挂载容易混淆，并且其传参方式限制了使用场景。

`--mount` 相对支持丰富的挂载方式

```sh
$ sudo docker run -d --name webapp webapp:latest --mount 'type=volume,src=appdata,dst=/webapp/storage,volume-driver=local,volume-opt=type=nfs,volume-opt=device=<nfs-server>:<nfs-path>' webapp:latest
```

在 `--mount` 中，我们可以通过逗号分隔这种 CSV 格式来定义多个参数。其中，通过 type 我们可以定义挂载类型，其值可以是：bind，volume 或 tmpfs。另外，`--mount` 选项能够帮助我们实现集群挂载的定义，例如在这个例子中，我们挂载的来源是一个 NFS 目录。

由于在实际开发中，`-v` 基本上足够满足我们的需求，所以我们不常使用相对复杂的 `--mount` 选项来定义挂载，这里我们只是将它简单介绍，供大家参考。



## 8.保存和共享镜像

Docker 镜像的本质是多个基于 UnionFS 的镜像层依次挂载的结果，而容器的文件系统则是在以只读方式挂载镜像后增加的一个可读可写的沙盒环境。

基于这样的结构，Docker 中为我们提供了将容器中的这个可读可写的沙盒环境持久化为一个镜像层的方法。

将容器修改的内容保存为镜像的命令是 `docker commit`，由于镜像的结构很像代码仓库里的修改记录，而记录容器修改的过程又像是在提交代码，所以这里我们更形象的称之为提交容器的更改。

```sh
$ sudo docker commit webapp
sha256:0bc42f7ff218029c6c4199ab5c75ab83aeaaed3b5c731f715a3e807dda61d19e
```

然后能在本地镜像列表找到它

```sh
$ sudo docker images
REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
<none>                <none>              0bc42f7ff218        3 seconds ago       372MB
## ......
```



还可以类似git一样，提交时给出一个提交信息

```shell
$ sudo docker commit -m "Configured" -a wuxie 容器ID 新名字
```

-a：提交镜像作者

### 为镜像命名

上方可以看到新持久化的镜像没有`REPOSITORY和TAG `，可以对其进行命名`docker tag`

```sh
$ sudo docker tag 0bc42f7ff218 webapp:1.0
```

也能更改已有的镜像名

```sh
$ sudo docker tag webapp:1.0 webapp:latest
```

创建新的命名，原本的旧的镜像依然存在列表中。

```sh
$ sudo docker images
REPOSITORY            TAG                 IMAGE ID            CREATED             SIZE
webapp                1.0                 0bc42f7ff218        29 minutes ago      372MB
webapp                latest              0bc42f7ff218        29 minutes ago      372MB
## ......
```

由于镜像是对镜像层的引用记录，所以我们对镜像进行命名后，虽然能够在镜像列表里同时看到新老两个镜像，实质是它们其实引用着相同的镜像层，看ID即可得知。

还可以在提交时进行命名

```sh
$ sudo docker commit -m "Upgrade" webapp webapp：2.0
```



### 镜像的迁移

---

镜像是由Docker管理，先从docker中将镜像去除`docker save`,保存到docker外部。

```sh
$ sudo docker save webapp:1.0 > webapp-1.0.tar
```

`docker save`会将镜像放在输出流中，需要用管道进行接收（`>`)，还提供`-o`(output)指定输出镜像。

```sh
$ sudo docker save -o ./webapp-1.0.tar webapp:1.0
```



#### 导入镜像

可以将镜像文件复制到另一台机器上，如此又要将镜像导入到这台机器的dockers中。

docker提供`dockre load` 将镜像导入docker中

```sh
$ sudo docker load < webapp-1.0.tar
```

同样也是将镜像放入输入流中，用`<` 读获得，还提供`-i(input)`指定选项输入文件。

```sh
$ sudo docker load -i webapp-1.0.tar
```



#### 批量迁移

可以将多个镜像打包成一个，便于一次性迁移多个镜像

```sh
$ sudo docker save -o ./images.tar webapp:1.0 nginx:1.12 mysql:5.7
```

`docker load` 可以识别并导入



### 导出和导入容器

---

 导出容器再导入，开发者觉得效率有些低，所以提供了对容器的直接导入和导出。

`docker export` 可以理解为`docker commit`和`docker save`的结合

```sh
$ sudo docker export -o ./webapp.tar webapp
```

`docker import` 并且直接将容器导入，而是将容器运行时的内容以镜像的形式导入。所以导入的结果其实是一个镜像。

```sh
$ sudo docker import ./webapp.tar webapp:1.0
```

> docker export 比 docker save 保存的包要小，原因是 save 保存的是一个分层的文件系统，export 导出的只是一层文件系统。

![img](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/1404523-20220529225546290-1147868208.png)

Docker镜像实际上就是由这样的一层层文件进行叠加起来的，上层的文件会覆盖下层的同名文件。如果将`postgres-save.tar`中的各层文件合并到一起，基本就是`postgres-export.tar`的内容。由于`postgres-save.tar`里面的各层文件会存在很多重复的文件，这也解释了为什么`postgres-save.tar`会比`postgres-export.tar`大100多M。







## 9.镜像仓库

---

#### 安装运行

```sh
$ docker run -d -p 5000:5000 --restart=always --name registry registry
```

这将使用官方的 `registry` 镜像来启动私有仓库。默认情况下，仓库会被创建在容器的 `/var/lib/registry` 目录下。你可以通过 `-v` 参数来将镜像文件存放在本地的指定路径。例如下面的例子将上传的镜像放到本地的 `/opt/data/registry` 目录。

```sh
$ docker run -d \
    -p 5000:5000 \
    -v /opt/data/registry:/var/lib/registry \
    registry
```



#### 上传，搜索，下载

使用`docker tag` 标记镜像，然后推送仓库。

例如私有仓库地址为 `127.0.0.1:5000`。使用 `docker tag` 将 `xx:latest` 这个镜像标记为 `127.0.0.1:5000/xxx:latest`。

```sh
$ docker tag ubuntu:latest 127.0.0.1:5000/ubuntu:latest
$ docker image ls
REPOSITORY                        TAG                 IMAGE ID            CREATED             VIRTUAL SIZE
ubuntu                            latest              ba5877dc9bec        6 weeks ago         192.7 MB
127.0.0.1:5000/ubuntu:latest      latest              ba5877dc9bec        6 weeks ago         192.7 MB
```



`docker push`上传标记的镜像

```sh
$ docker push 127.0.0.1:5000/ubuntu:latest
The push refers to repository [127.0.0.1:5000/ubuntu]
373a30c24545: Pushed
a9148f5200b0: Pushed
cdd3de0940ab: Pushed
fc56279bbb33: Pushed
b38367233d37: Pushed
2aebd096e0e2: Pushed
latest: digest: sha256:fe4277621f10b5026266932ddf760f5a756d2facd505a94d2da12f4f52f71f5a size: 1568
```



配置非https仓库地址，因为 Docker 默认不允许非 `HTTPS` 方式推送镜像；

> 对于使用 `systemd` 的系统，请在 `/etc/docker/daemon.json` 中写入如下内容
>
> windwos在docker desktop设置里面Docker Engine 即可修改

```sh
{
  "registry-mirrors": [
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ],
  "insecure-registries": [
    "192.168.199.100:5000"
  ]
}
```



查看仓库

```sh
$ curl 127.0.0.1:5000/v2/_catalog

StatusCode        : 200
StatusDescription : OK
Content           : {"repositories":["redis"]}

RawContent        : HTTP/1.1 200 OK
                    Docker-Distribution-Api-Version: registry/2.0
                    X-Content-Type-Options: nosniff
                    Content-Length: 27
                    Content-Type: application/json; charset=utf-8
                    Date: Thu, 16 Mar 2023 07:13:15 GMT...
Forms             : {}
Headers           : {[Docker-Distribution-Api-Version, registry/2.0], [X-Content-Type-Options, nosnif
                    f], [Content-Length, 27], [Content-Type, application/json; charset=utf-8]...}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 27
```



# 高级篇

## 1.MySQL主从复制

---

新建主服务器实例

```sh
docker run -p 3307:3306 --name mysql-master -v /mydata/mysql-master/log:/var/log/mysql -v /mydata/mysql-master/data:/var/lib/mysql -v /mydata/mysql-master/conf:/etc/mysql -e MYSQL_ROOT_PASSWORD=root  -d mysql:5.7
```

进入`/mydata/mysql-master/conf` 创建my.cnf

```sh
[mysqld]
## 设置server_id，同一局域网中需要唯一
server_id=101 
## 指定不需要同步的数据库名称
binlog-ignore-db=mysql  
## 开启二进制日志功能
log-bin=mall-mysql-bin  
## 设置二进制日志使用内存大小（事务）
binlog_cache_size=1M  
## 设置使用的二进制日志格式（mixed,statement,row）
binlog_format=mixed  
## 二进制日志过期清理时间。默认值为0，表示不自动清理。
expire_logs_days=7  
## 跳过主从复制中遇到的所有错误或指定类型的错误，避免slave端复制中断。
## 如：1062错误是指一些主键重复，1032错误是因为主从数据库数据不一致
slave_skip_errors=1062
```

`docker restart mysql-master` 重启

进入容器，创建数据同步用户

```sh
CREATE USER 'slave'@'%' IDENTIFIED BY '123456';
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'slave'@'%';
```



创建从服务器

```sh
docker run -p 3308:3306 --name mysql-slave -v /mydata/mysql-slave/log:/var/log/mysql -v /mydata/mysql-slave/data:/var/lib/mysql -v /mydata/mysql-slave/conf:/etc/mysql -e MYSQL_ROOT_PASSWORD=root  -d mysql:5.7
```

配置my.cnf

```sh
[mysqld]
## 设置server_id，同一局域网中需要唯一
server_id=102
## 指定不需要同步的数据库名称
binlog-ignore-db=mysql
## 开启二进制日志功能，以备Slave作为其它数据库实例的Master时使用
log-bin=mall-mysql-slave1-bin
## 设置二进制日志使用内存大小（事务）
binlog_cache_size=1M
## 设置使用的二进制日志格式（mixed,statement,row）
binlog_format=mixed
## 二进制日志过期清理时间。默认值为0，表示不自动清理。
expire_logs_days=7
## 跳过主从复制中遇到的所有错误或指定类型的错误，避免slave端复制中断。
## 如：1062错误是指一些主键重复，1032错误是因为主从数据库数据不一致
slave_skip_errors=1062
## relay_log配置中继日志
relay_log=mall-mysql-relay-bin
## log_slave_updates表示slave将复制事件写进自己的二进制日志
log_slave_updates=1
## slave设置为只读（具有super权限的用户除外）
read_only=1

```



`show master status` 主服务器查看主从同步状态

`show slave status` 从服务器查看主从同步状态；

从服务器配置主从复制

```sh
master_host：主数据库的IP地址；
master_port：主数据库的运行端口；
master_user：在主数据库创建的用于同步数据的用户账号；
master_password：在主数据库创建的用于同步数据的用户密码；
master_log_file：指定从数据库要复制数据的日志文件，通过查看主数据的状态，获取File参数；
master_log_pos：指定从数据库从哪个位置开始复制数据，通过查看主数据的状态，获取Position参数；
master_connect_retry：连接失败重试的时间间隔，单位为秒。 
```

```sh
change master to master_host='192.168.88.188', master_port=3306,master_user='slave',master_password='MySQL5.7clone',master_log_file='mysql-bin.000001',master_log_pos=617,master_connect_retry=30;

```



![image-20230317095256323](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317095256323.png)

这里表示开启了。

主服务器建表插入数据，从服务器可以查看到。

## 2.redis集群

> 1~2亿条数据需要缓存，请问如何设计这个存储案例
>
> 肯定使用分布式缓存，redis如何落地？

### 1.存储数据

#### 方案一（哈希取余分区)

---

`hash(key)%n` n节点数，对redis的key进行哈希值的取余来判断该存储在哪一个redis节点上。

缺点：

 原来规划好的节点，进行扩容或者缩容就比较麻烦了额，不管扩缩，每次数据变动导致节点有变动，映射关系需要重新进行计算，在服务器个数固定不变时没有问题，**如果需要弹性扩容或故障停机的情况下，原来的取模公式就会发生变化**：Hash(key)/3会变成Hash(key) /?。此时地址经过取余运算的结果将发生很大变化，根据公式获取的服务器也会变得不可控。某个redis机器宕机了，由于台数数量变化，会导致hash取余全部数据重新洗牌。 

#### 方案二（一致性哈希算法区分）

也是按照取模的方法，但方案一是对节点的数量进行取模。而一致性哈希算法是对2^32的取模，按照常用的hash算法来将对应的key哈希到一个具有**2^32次方个桶**的空间中，即0~(2^32)-1的数字空间中。现在我们可以将这些数字头尾相连，想象成一个闭合的环形。如下图：

![img](https://pic1.zhimg.com/80/v2-ae2d4482fe02c5a631797e264e3ada24_720w.webp)

> 解释为什么是**2^32** ，因为是对服务器进行负载均衡，服务器IP地址是32位，所以是2^32-1的数值空间。

**key落键规则**

当我们需要存储一个kv键值对时，首先计算key的hash值，hash(key)，将这个key使用相同的函数Hash计算出哈希值并确定此数据在环上的位置，从此位置沿环**顺时针**“行走”，**第一台遇到的服务器就是其应该定位到的服务器**，并将该键值对存储在该节点上。

![img](https://pic4.zhimg.com/80/v2-50d10655c9faa9dca36fb4fd32ca15bf_720w.webp)

**优点**

- 容错性：如上图，如果t2宕机了，那么m3,m4就会去寻找t2下一台机器t1,只会影响到**t3到t2这段距离的映射**，而其他不会影响。
- 扩容性：扩展性数据量增加了，需要增加一台节点NodeX，X的位置在A和B之间，那受到影响的也就是A到X之间的数据，重新把A到X的数据录入到X上即可，不会导致hash取余全部数据重新洗牌。

**缺点**

- 数据倾斜问题：节点太小，导致数据集中存储在某一个节点上，导致另一台节点数据太少。 

![img](https://pic1.zhimg.com/80/v2-ceaa7afa8ba1e5dfeec1bed221bb49f8_720w.webp)



#### 方案三（哈希槽分区）

##### 1.介绍

出现背景：解决一致性哈希算法的数据倾斜问题；

本质是一个数组，**数组[0,2^14-1]形成的hash slot空间**



##### 2.作用

解决均匀分配问题，在数据和节点之间又加入了一层，称为哈希槽用于管理数据和节点之间的关系。相当于节点上放的是槽，槽里放的是数据。

槽解决粒度问题，将原本节点的数据粒度变大，便于数据移动。

哈希解决的是映射问题，使用key的哈希值计算所在槽，便于进行数据的分配。



##### 3.多少个hash槽

> 为什么是2^14

在redis节点发送心跳包时需要把所有的槽放到这个心跳包里，以便让节点知道当前集群信息，16384=16k，在发送心跳包时使用`char`进行bitmap压缩后是2k（`2 * 8 (8 bit) * 1024(1k) = 16K`），也就是说使用2k的空间创建了16k的槽数。

虽然使用CRC16算法最多可以分配65535（2^16-1）个槽位，65535=65k，压缩后就是8k（`8 * 8 (8 bit) * 1024(1k) =65K`），也就是说需要需要8k的心跳包，作者认为这样做不太值得；并且一般情况下一个redis集群不会有超过1000个master节点，所以16k的槽位是个比较合适的选择。



哈希槽计算

CRC(key)%16384

![image-20230317164253507](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317164253507.png)

 

#### 实践

##### 数据读写存储

1. **新建6个redis实例**

```sh
docker run -d --name redis-node-1 --net host --privileged=true -v /data/redis/share/redis-node-1:/data redis --cluster-enabled yes --appendonly yes --port 6381

docker run -d --name redis-node-2 --net host --privileged=true -v /data/redis/share/redis-node-2:/data redis --cluster-enabled yes --appendonly yes --port 6382
docker run -d --name redis-node-3 --net host --privileged=true -v /data/redis/share/redis-node-3:/data redis --cluster-enabled yes --appendonly yes --port 6383
 
docker run -d --name redis-node-4 --net host --privileged=true -v /data/redis/share/redis-node-4:/data redis --cluster-enabled yes --appendonly yes --port 6384
 
docker run -d --name redis-node-5 --net host --privileged=true -v /data/redis/share/redis-node-5:/data redis --cluster-enabled yes --appendonly yes --port 6385
 
docker run -d --name redis-node-6 --net host --privileged=true -v /data/redis/share/redis-node-6:/data redis --cluster-enabled yes --appendonly yes --port 6386

```

命令解释：

- `--net host` ·使用宿主机的IP和端口，默认

- `--privileged=true` 获得宿主机的root权限
- `--cluster-enabled yes` 开启redis集群
- `--appendonly yes` 持久化



2. **构建主从关系**

进入一个实例,执行命令

```sh
redis-cli --cluster create ip:6381 ip:6382 ip:6383 ip:6384 ip:6385 ip:6386 --cluster-replicas 1
```

`--cluster-replicas 1 表示为每个master创建一个slave节点`

![image-20230317171447570](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317171447570.png)



进入节点一查看集群状态

```sh
$ redis-cli -p 6381
127.0.0.1:6381> cluster info
127.0.0.1:6381> cluster nodes 查看节点状态
```

![image-20230317175022691](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317175022691.png)

自动分配挂载点

3.**开启集群环境**

```sh
127.0.0.1:6381> set k1 v1
(error) MOVED 12706 ip:6383
127.0.0.1:6381> set k2 v1
OK
```

因为我们配置了集群，redis哈希槽也进行了分配，那么有些内容就存不近当前节点。所以使用参数`-c` 进入集群反击,防止路由失效;

```sh
root@localhost:/data# redis-cli -p 6381 -c
127.0.0.1:6381> set k1 v1
-> Redirected to slot [12706] located at ip:6383
OK
ip:6383> keys *
1) "k1"
ip:6383> 
```

4.**查看集群信息**

```sh
root@localhost:/data# redis-cli --cluster check ip:6381
ip:6381 (0c472300...) -> 0 keys | 5461 slots | 1 slaves.
ip:6383 (de42cb3d...) -> 1 keys | 5461 slots | 1 slaves.
ip:6382 (e2eb9dc3...) -> 0 keys | 5462 slots | 1 slaves.
```



##### 容错切换迁移

**问题：父节点宕机后，子节点是否能上位？**

停掉node1，登录node2，输入查看集群信息指令

![image-20230317204524134](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317204524134.png)

可见子节点上位称为父节点。



**问题：如果6381节点重新启动，那么集群关系会如何？**

![image-20230317205602289](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317205602289.png)

可见不会恢复到原来的父节点关系



### 2.主从扩容

1. **新增主机**

```sh
docker run -d --name redis-node-7 --net host --privileged=true -v /data/redis/share/redis-node-7:/data redis --cluster-enabled yes --appendonly yes --port 6387

docker run -d --name redis-node-8 --net host --privileged=true -v /data/redis/share/redis-node-8:/data redis --cluster-enabled yes --appendonly yes --port 6388
```

将新增节点（空槽位）作为master加入原集群

```sh
redis-cli --cluster add-node 自己实际IP地址:6387 自己实际IP地址:6381
一开始是以6381作为集群的主节点，根据他找到对应的集群
```

查看集群情况,即可见到四个父节点，但是没有槽位

```sh
redis-cli --cluster check ip:6381
```

2.**重新分配槽号**

```sh
redis-cli --cluster reshard ip:6381
```

![image-20230317211553269](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317211553269.png)



检查集群信息

![image-20230317211644870](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317211644870.png)

**其他三个节点各自分配了一点槽位给新节点。**

> 因为重新分配的成本太高



3.**添加子节点**

```sh
redis-cli --cluster add-node ip:新slave端口 ip:新master端口 --cluster-slave --cluster-master-id 新主机节点ID
```



### 3.主从缩容

1. **先删除从节点**

```sh
redis-cli --cluster del-node ip:从机端口 从机6388节点ID
ip:6388 535f9ae790da5409dabd7ba7a45eacc6db020535
```

查看集群信息得S只有三个了

2. **父节点清空槽位并进行分配**

```sh
redis-cli --cluster reshard ip:6381
```

![image-20230317215109390](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317215109390.png)



![image-20230317215130505](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230317215130505.png)

然后删除父节点

```sh
redis-cli --cluster check ip:6381
```



## 3.使用Dockerfile 定制镜像

### 编写dockerfile

镜像的定制实际上就是定制每一层所添加的配置、文件。如果我们可以把每一层修改、安装、构建、操作的命令都写入一个脚本，用这个脚本来构建、定制镜像，那么之前提及的无法重复的问题、镜像构建透明性的问题、体积的问题就都会解决。这个脚本就是 Dockerfile。

`Dockerfile` 是一个**文本文件(没有任何后缀)**，包含一条条指令，每一条指令构建一层，因此每一条指令内容就是描述该层如何构建。

定制nginx镜像为例，使用Dockerfile定制

```sh
FROM nginx
RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
```

**解释指令**

1. `FROM` 指定基础镜像，定制镜像必须在某一个镜像基础之上，所以`FROM` 是必备指令，并且必须是第一条指令。除了选择现有镜像为基础镜像外，Docker 还存在一个特殊的镜像，名为 `scratch`。这个镜像是虚拟的概念，并不实际存在，它表示一个**空白的镜像**。

2. `RUN` 容器构建时需要运行的命令，其格式有两种

   1. shell格式，`RUN 命令`

   ```sh
   RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
   ```

   2. exec格式,`RUN ["可执行文件", "参数1", "参数2"]` ,类似函数调用的格式

3. 每条保留字指令必须为**大写字母**



```sh
FROM debian:stretch

RUN apt-get update
RUN apt-get install -y gcc libc6-dev make wget
RUN wget -O redis.tar.gz "http://download.redis.io/releases/redis-5.0.3.tar.gz"
RUN mkdir -p /usr/src/redis
RUN tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
RUN make -C /usr/src/redis
RUN make -C /usr/src/redis install
```

如上代码是没有意义的，因为前面说过`Dockerfile` 每一个指令都会建立一层，`RUN` 也一样，这样就建立了七层镜像，将需要运行时不需要的东西,比如编译环境、更新的软件包等等都装入了镜像，导致结果很臃肿，增加了构建时间，也容易出错。

> Union FS 是有最大层数限制的，比如 AUFS，曾经是最大不得超过 42 层，现在是不得超过 127 层。

上面代码正常编写应该如下

```sh
FROM debian:stretch

RUN set -x; buildDeps='gcc libc6-dev make wget' \
    && apt-get update \
    && apt-get install -y $buildDeps \
    && wget -O redis.tar.gz "http://download.redis.io/releases/redis-5.0.3.tar.gz" \
    && mkdir -p /usr/src/redis \
    && tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
    && make -C /usr/src/redis \
    && make -C /usr/src/redis install \
    && rm -rf /var/lib/apt/lists/* \
    && rm redis.tar.gz \
    && rm -r /usr/src/redis \
    && apt-get purge -y --auto-remove $buildDeps
```

将具有相同目的的内容编写在一起，如上是为了编译安装redis可执行文件，就没必要编写多层，将其合为一个`RUN`。提醒自己`Dockerfile` 并非是编写shell脚本，而是**定义每一层如何构建**。

也注意编写格式，`Dockerfile` 提供`\,#` 等符号进行格式的调整。



```sh
    && rm -rf /var/lib/apt/lists/* \
    && rm redis.tar.gz \
    && rm -r /usr/src/redis \
    && apt-get purge -y --auto-remove $buildDeps
```



此外，还可以看到这一组命令的最后添加了清理工作的命令，删除了为了编译构建所需要的软件，清理了所有下载、展开的文件，并且还清理了 `apt` 缓存文件。这是很重要的一步，我们之前说过，镜像是多层存储，**每一层的东西并不会在下一层被删除**，会一直跟随着镜像。因此镜像构建时，**一定要确保每一层只添加真正需要添加的东西**，任何无关的东西都应该清理掉。



#### 相关命令

```sh
1.FROM
2.RUN
3.EXPOSE 暴露端口
4.WORKDIR  指定在创建容器后，终端默认登陆的进来工作目录，一个落脚点
5.USER 指定该镜像以什么样的用户去执行，如果都不指定，默认是root
6.ENV 用来在构建镜像过程中设置环境变量
ENV MY_PATH /usr/mytest
这个环境变量可以在后续的任何RUN指令中使用，这就如同在命令前面指定了环境变量前缀一样；
也可以在其它指令中直接使用这些环境变量，
比如：WORKDIR $MY_PATH

7.VOLUME 容器数据卷
8.ADD/COPY
ADD-将宿主机目录下的文件拷贝进镜像且会自动处理URL和解压tar压缩包(COPY+解压)
COPY-将文件进行拷贝进镜像 COPY src dest/COPY ["src", "dest"]

9.CMD 指定容器启动后要做的事
Dockerfile可以有多个CMD指令，但只有最后一个生效，CMD 会被 docker run 之后的参数替换

Tomcat dockerfile
EXPOSE 8080
CMD ["catalina.sh", "run"]
执行tomcat docker run -it -p 8080:8080 tomcat /bin/bash
如此就导致最后启动tomcat的CMD被/bin/bash覆盖无法正常访问tomcat

最后注意CMD是docker run时运行，RUN是docker build时运行

10.ENTRYPOINT 也是用来指定一个容器启动时要运行的命令
类似于 CMD 指令，但是ENTRYPOINT不会被docker run后面的命令覆盖，
而且这些命令行参数会被当作参数送给 ENTRYPOINT 指令指定的程序

11.MAINTAINER 作者信息 
```

![image-20230318090454014](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318090454014.png)



| 是否传参         | 按照dockerfile编写执行         | 传参运行                                      |
| ---------------- | ------------------------------ | --------------------------------------------- |
| Docker命令       | docker run  nginx:test         | docker run  nginx:test -c /etc/nginx/new.conf |
| 衍生出的实际命令 | nginx -c /etc/nginx/nginx.conf | nginx -c /etc/nginx/new.conf                  |



#### 自定义centos:seven: 自带java8，ifconfig，vim

```sh
FROM centos
MAINTAINER zzyy<zzyybs@126.com>
 
ENV MYPATH /usr/local
WORKDIR $MYPATH
 
#安装vim编辑器
RUN yum -y install vim
#安装ifconfig命令查看网络IP
RUN yum -y install net-tools
#安装java8及lib库

RUN yum -y install glibc.i686
RUN mkdir /usr/local/java
#ADD 是相对路径jar,把jdk-8u171-linux-x64.tar.gz添加到容器中,安装包必须要和Dockerfile文件在同一位置
ADD jdk-8u171-linux-x64.tar.gz /usr/local/java/
#配置java环境变量
ENV JAVA_HOME /usr/local/java/jdk1.8.0_171
ENV JRE_HOME $JAVA_HOME/jre
ENV CLASSPATH $JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar:$JRE_HOME/lib:$CLASSPATH
ENV PATH $JAVA_HOME/bin:$PATH
 
EXPOSE 80
 
CMD echo $MYPATH
CMD echo "success--------------ok"
CMD /bin/bash
```



#### 虚悬镜像

> 在镜像构建时出现一些错误，导致仓库名、标签都是<none>的镜像，俗称dangling 

```sh
from ubuntu
CMD echo 'action is success'
```

![image-20230318095245261](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318095245261.png)

**查看**

```sh
docker image ls -f dangling=true
```

**删除**

```sh
docker image prune
```



### 构建镜像

```sh
$ docker build -t nginx:v3 .
Sending build context to Docker daemon 2.048 kB
Step 1 : FROM nginx
 ---> e43d811ce2f4
Step 2 : RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
#启动了一个容器进行执行命令
 ---> Running in 9cdc27646c7b
 ---> 44aa4490ce2c
 #删除了容器
Removing intermediate container 9cdc27646c7b
Successfully built 44aa4490ce2c
```

![image-20230318084244081](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318084244081.png)



```sh
docker build [选项] <上下文路径/URL/-> -t 名字
```

#### docker build 的工作原理

1. Client端执行 docker build . 命令 ;
2. Docker 客户端会将构建命令后面指定的路径(.)下的所有文件打包发送给 Docker 服务端;
3. Docker 服务端收到客户端发送的包，然后解压，根据 Dockerfile 里面的指令进行镜像的分层构建；

![https://docs.docker.com/engine/images/architecture.svg](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/architecture.svg)

### 镜像构建上下文

docker构建最后会指定 `.` 。虽然 `.` 表示当前目录，但实际这个路径并不是`Dockerfile` 的所在路径，而是**上下文路径** 。 

上下文路径，是指 `docker` 在构建镜像，有时候想要使用到本机的文件（比如复制），docker build 命令得知这个路径后，会将**路径下的所有内容打包**。

**解析：**由于 docker 的运行模式是 C/S。我们本机是 C（client），docker 引擎是 S（server）。实际的构建过程是在 docker 引擎下完成的，所以这个时候无法用到我们本机的文件。这就需要把我们本机的指定目录下的文件一起打包提供给 docker 引擎使用。

如果未说明最后一个参数，那么**默认上下文路径就是 Dockerfile 所在的位置**。

**注意：**上下文路径下不要放无用的文件，因为会一起打包发送给 docker 引擎，如果文件过多会造成过程缓慢。



因为我们编写Dockerfile会用到一些`COPY,ADD` 命令，将本地文件复制进镜像，而此时构建并非在本地，所以就需要用到上下文。

```shell
COPY ./package.json /app/
```

这并不是要复制执行 `docker build` 命令所在的目录下的 `package.json`，也不是复制 `Dockerfile` 所在目录下的 `package.json`，而是复制 **上下文（context）** 目录下的 `package.json`。

如果希望目录下有些东西不希望打包，可以使用`.gitignore` 一样的语法写一个 `.dockerignore`，该文件是用于剔除不需要作为上下文传递给 Docker 引擎的。

> 实际上 `Dockerfile` 的文件名并不要求必须为 `Dockerfile`，而且并不要求必须位于上下文目录中，比如可以用 `-f ../Dockerfile.php` 参数指定某个文件作为 `Dockerfile`。



其他`docker build` 用法

1. `docker build` 还支持从 URL 构建，比如可以直接从 Git repo 中构建：

```sh
# $env:DOCKER_BUILDKIT=0
# export DOCKER_BUILDKIT=0

$ docker build -t hello-world https://github.com/docker-library/hello-world.git#master:amd64/hello-world

Step 1/3 : FROM scratch
 --->
Step 2/3 : COPY hello /
 ---> ac779757d46e
Step 3/3 : CMD ["/hello"]
 ---> Running in d2a513a760ed
Removing intermediate container d2a513a760ed
 ---> 038ad4142d2b
Successfully built 038ad4142d2b
```

`#master` 指定分支，构建目录为`amd64/hello-world`。然后docker就自动git clone

2. 用给定tar压缩包构建

```sh
$ docker build http://server/context.tar.gz
```

3. 从标准输入中读取 Dockerfile 进行构建

```sh
docker build - < Dockerfile

===
cat Dockerfile | docker build -
```

`-` 占位

4. 从标准输入中读取上下文压缩包进行构建

```sh
$ docker build - < context.tar.gz
```



### 实践 部署springboot项目

![image-20230318105945447](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318105945447.png)

```sh
# 基础镜像使用java
FROM java:8
# 作者
MAINTAINER wuxie
# VOLUME 指定临时文件目录为/tmp，在主机/var/lib/docker目录下创建了一个临时文件并链接到容器的/tmp
VOLUME /tmp
# 将jar包添加到容器中并更名为xx.jar
ADD E-commence-center.jar wuxie_docker.jar
# 运行jar包
# 将本机文件复制到docker内部
RUN bash -c 'touch /wuxie_docker.jar'
ENTRYPOINT ["java","-jar","/wuxie_docker.jar"]
#暴露6001端口作为微服务
EXPOSE 6001
 
```



## 4.Docker网络

### 配置网络

#### 容器网络

---

容器网络实质上也是由 Docker 为应用程序所创造的**虚拟环境的一部分**，它能让应用从宿主机操作系统的网络环境中独立出来，形成容器自有的网络设备、IP 协议栈、端口套接字、IP 路由表、防火墙等等与网络相关的模块。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/5/165a810ad2c81714~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

Docker网络中，三个核心概念:**沙盒 ( Sandbox )**、**网络 ( Network )**、**端点 ( Endpoint )**

- **沙盒**提供了容器的虚拟网络栈，也就是之前所提到的端口套接字、IP 路由表、防火墙等的内容。其实现隔离了容器网络与宿主机网络，形成了完全独立的容器网络环境。
- **网络**可以理解为 Docker **内部的虚拟子网**，网络内的参与者相互可见并能够进行通讯。Docker 的这种虚拟网络也是于宿主机网络存在隔离关系的，其目的主要是形成容器间的安全通讯环境。
- **端点（端口)**是位于容器或网络隔离墙之上的洞，其主要目的是形成一个可以控制的突破封闭的网络环境的出入口。当容器的端点与网络的端点形成配对后，就如同在这两者之间搭建了桥梁，便能够进行数据传输了。

#### 网络实现

容器网络模型为容器引擎提供了一套标准的网络对接范式，而在 Docker 中，实现这套范式的是 Docker 所封装的 **libnetwork 模块**。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/23/166042a49627f8a6~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

目前 Docker 官方为我们提供了五种 Docker 网络驱动，分别是：**Bridge Driver**、**Host Driver**、**Overlay Driver**、**MacLan Driver**、**None Driver**。

- Bridge 网络是 Docker 容器的**默认网络驱动**，简而言之其就是通过网桥来实现网络通讯 ( 网桥网络的实现可以基于硬件，也可以基于软件 )。
- Overlay 网络是借助 Docker 集群模块 Docker Swarm 来搭建的跨 Docker Daemon 网络，我们可以通过它搭建跨物理主机的虚拟网络，进而让不同物理机中运行的容器感知不到多个物理机的存在。



### 容器互联

---

由于 Docker 提倡容器与应用共生的轻量级容器理念，所以容器中通常只包含一种应用程序。但目前系统服务而言，不可能靠单一应用支撑，而是多个应用相互组成的。所以在Docker中需要通过多个容器组成这样的系统。

`docker create` 或 `docker run` 创建时通过 `--link` 选项进行配置。

```sh
$ sudo docker run -d --name mysql -e MYSQL_RANDOM_ROOT_PASSWORD=yes mysql
$ sudo docker run -d --name webapp --link mysql webapp:latest
```

Docker提供友好的连接方式，只需要将容器的网络命名填入连接地址即可。

```java
String url = "jdbc:mysql://mysql:3306/webapp";
```

通过`mysql:3306`,docker自己解析并指向MySQL容器的IP地址**。如此在我们更换环境时，我们就不需要进行更改IP地址，只需要别名，然后Docker就自动映射**。



#### 相关命令

```sh
docker network ls
docker network inspect xx名字 查看详细信息
docker network rm xx
```



#### 暴露端口

知道了IP地址，不意味就可以连接了。如我们的电脑有防火墙，docker为容器网络也增加了一套安全机制，类似防火墙。我们需要暴露容器的端口，才可以使得其他容器访问。

`docker ps`即可看到

```sh
$docker ps -a
CONTAINER ID   IMAGE          COMMAND                  CREATED        STATUS         PORTS                     NAMES
610dff2378de   redis:latest   "docker-entrypoint.s…"   38 hours ago   Up 2 seconds   0.0.0.0:32769->6379/tcp   redis-80wo
```

如上`0.0.0.0:32769->6379/tcp`,暴露了32769和6379。连接时就可以对这两个端口访问。

也可以在容器创建时定义`--expose`

```sh
$ sudo docker run -d --name mysql -e MYSQL_RANDOM_ROOT_PASSWORD=yes --expose 13306 --expose 23306 mysql:5.7
```



#### 别名连接

```sh
$ sudo docker run -d --name webapp --link mysql:database webapp:latest
```

`mysql:databse`，用户可自定义别名，然后进行连接。较为灵活

```java
String url = "jdbc:mysql://database:3306/webapp";
```



#### 管理网络

网络这个概念我们可以理解为 Docker 所虚拟的子网，而容器网络沙盒可以看做是虚拟的主机，只有当**多个主机在同一子网里时，才能互相看到并进行网络数据交换**。

当我们启动 Docker 服务时，它会为我们**创建一个默认的 bridge 网络，而我们创建的容器在不专门指定网络的情况下都会连接到这个网络上**。

`docker inspect` Network部分看到容器网络相关信息

```sh
$ sudo docker inspect mysql
[
    {
## ......
        "NetworkSettings": {
## ......
            "Networks": {
                "bridge": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": null,
                    "NetworkID": "bc14eb1da66b67c7d155d6c78cb5389d4ffa6c719c8be3280628b7b54617441b",
                    "EndpointID": "1e201db6858341d326be4510971b2f81f0f85ebd09b9b168e1df61bab18a6f22",
                    "Gateway": "172.17.0.1",
                    "IPAddress": "172.17.0.2",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "02:42:ac:11:00:02",
                    "DriverOpts": null
                }
            }
## ......
        }
## ......
    }
]
```



#### 网络分配规则

创建两个实例查看ip地址

![image-20230318113359263](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318113359263.png)

![image-20230318113425648](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318113425648.png)

关闭u2后再新建一个实例查看ip变化，**可以看到docker内部网络会发现改变。那么有些服务就可能无效。所以就需要自己规划好网络**

![image-20230318113433298](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318113433298.png)



#### Docker相关网络

##### 1.bridge

> 1.Docker 服务默认会创建一个 docker0 网桥（其上有一个 docker0 内部接口），该桥接网络的名称为docker0，它在内核层连通了其他的物理或虚拟网卡，这就将所有容器和本地主机都放到同一个物理网络。Docker 默认指定了 docker0 接口 的 IP 地址和子网掩码，让主机和容器之间可以通过网桥相互通信。
>
> 2.Docker启动一个容器时会根据Docker网桥的网段分配给容器一个IP地址，称为Container-IP，同时Docker网桥是每个容器的默认网关。因为在同一宿主机内的容器都接入同一个网桥，这样容器之间就能够通过容器的Container-IP直接通信。
>
> **docker run** 不指定network默认就是bridge，就是docker0
>
> 3.网桥docker0创建一对对等虚拟设备接口一个叫veth，另一个叫eth0，成对匹配。这样一对接口叫`veth pair`

![image-20230318115356795](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318115356795.png)

测试

启动两个tomcat

```sh
docker run -d -p 8081:8080   --name tomcat81 billygoo/tomcat8-jdk8
docker run -d -p 8082:8080   --name tomcat82 billygoo/tomcat8-jdk8
```



![image-20230318121746508](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318121746508.png)

可以看到veth和eth0配对。

```sh
宿主机上:
# 28:veth14fc73a@if27
容器中:
# 27:eth0@if28
```





##### 2.host

**介绍**

> 直接使用宿主机的 IP 地址与外界进行通信，不再需要额外进行NAT 转换。
>
> 容器将不会获得一个独立的Network Namespace， 而是和宿主机共用一个Network Namespace。容器将不会虚拟出自己的网卡而是使用宿主机的IP和端口



![image-20230319091534575](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319091534575.png)



**使用**

1. `docker run -d -p 8083:8080 --network host --name tomcat83 billygoo/tomcat8-jdk8` 使用会出现警告但不会报错，因为如果使用host网络模式，那么端口号主要以主机端口为主，-p设置的就不会起作用，端口重复时会进行递增处理。
2. 正确：`docker run -d --network host --name tomcat83 billygoo/tomcat8-jdk8`

那么就可以直接在主机上访问`localhost:8080` tomcat界面

**windows的Docker for Windows，会自己配置端口映射。**



##### **3.none**

不做任何网络配置，只有一个`io`,通过`127.0.0.1:xx`访问



##### **4.container**

> 新建的容器**和已经存在的一个容器共享一个网络ip配置**而不是和宿主机共享。新创建的容器不会创建自己的网卡，配置自己的IP，而是和一个指定的容器共享IP、端口范围等。同样，两个容器除了网络方面，其他的如文件系统、进程列表等还是隔离的。**注意会出现端口冲突**

![image-20230319093803070](C:/Users/wuxie/AppData/Roaming/Typora/typora-user-images/image-20230319093803070.png)



**使用**

```sh
Alpine操作系统是一个面向安全的轻型Linux发行版 5.59MB
docker run -it   --name alpine1  alpine /bin/sh
docker run -it --network container:alpine1 --name alpine2  alpine /bin/sh

```

![image-20230319094532149](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319094532149.png)



**关闭alpine1查看alpine2的网络情况，可见网络配置消息**

![image-20230319094624593](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319094624593.png)

#### 自定义网络

##### 1.link已过时

![image-20230319095348158](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319095348158.png)



##### 2.出现背景

> ​	每个服务的ip会出现波动，那么理想情况就是对服务名进行调用，并不关注底层的ip地址，所以就出现自定义网络。

```sh
# 1. 新建
docker network create xxx
# 2. 启动
docker run -d -p 8081:8080 --network zzyy_network  --name tomcat81 billygoo/tomcat8-jdk8
docker run -d -p 8082:8080 --network zzyy_network  --name tomcat81 billygoo/tomcat8-jdk8
# 3.测试
ping tomcat81 成功
```

**自定义网络本身就维护好了主机名和ip的对应关系（ip和域名都能通）**



#### 创建网络

`docker network create` 

![image-20230318112641531](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230318112641531.png)

```sh
$ sudo docker network create -d bridge individual
```

`-d`为新的网络指定驱动，其值就上面提到的bridge等 默认为bridge

`docker network ls/list(list Windows用不了)`

```sh
$ docker network ls
NETWORK ID     NAME         DRIVER    SCOPE
4672a6c23d6a   bridge       bridge    local
400cb3365f9e   host         host      local
17d8187ce8d2   individual   bridge    local
87a00a1f1024   none         null      local
```

创建容器也可以指定加入的网络

```sh
$ sudo docker run -d --name mysql -e MYSQL_RANDOM_ROOT_PASSWORD=yes --network individual mysql:5.7
```

然后再通过`--link`让处于另一个网络的容器连接到该容器上，显示是失败的

```sh
$ sudo docker run -d --name webapp --link mysql --network bridge webapp:latest
docker: Error response from daemon: Cannot link to /mysql, as it does not belong to the default network.
ERRO[0000] error waiting for container: context canceled
```



### 端口映射

---

刚刚都是容器之间的网络访问，我们一般常用即是容器外通过网络访问容器内的应用。docker就提供了端口映射来实现。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/23/16605128077de72a~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

可以将容器的端口映射到宿主操作系统端口，从外部访问宿主操作系统的端口，数据请求就会自动放松给相关联的容器端口。

创建容器时`-p/-publish`进行映射

`-P`大写p是随机端口映射

```sh
$ sudo docker run -d --name nginx -p 80:80 -p 443:443 nginx:1.12
```

使用端口映射选项的格式是 `-p <ip>:<host-port>:<container-port>`，其中 ip 是宿主操作系统的监听 ip，可以用来控制监听的网卡，**默认为 0.0.0.0**，监听所有网卡。host-port 和 container-port 分别表示映射到宿主操作系统的端口和容器的端口，这两者是可以不一样的，我们可以将容器的 80 端口映射到宿主操作系统的 8080 端口，传入 `-p 8080:80` 即可。



```sh
$ sudo docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                                      NAMES
bc79fc5d42a6        nginx:1.12          "nginx -g 'daemon of…"   4 seconds ago       Up 2 seconds        0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp   nginx
```

`0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp   nginx`



#### windows和macos中使用映射

上方映射是在Linux系统的端口。而Windows和macos上运行的docker，Linux环境是虚拟出来，所以映射端口也只是虚拟环境中，并不能直接通过windows端口进行访问。

解决方案就是再加一次映射，将虚拟环境Linux系统的端口再映射到Windows和macos上。

![img](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2018/9/23/166053965573b1f4~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

Docker for Windows 会自动进行这步操作，所以不用担心。



## 5.Docker-compose 容器编排

---



### 1.概念

> Docker-Compose是Docker官方的开源项目，负责实现对Docker容器集群的快速编排。
>
> 可以管理多个 Docker 容器组成一个应用。你需要定义一个 YAML 格式的配置文件docker-compose.yml，**写好多个容器之间的调用关系**。然后，只要一个命令，就能同时启动/关闭这些容器

​	在执行一个服务时，你要先启动他所关联的其他服务器如MySQL，redis等，才能正常启动服务。所以compose就是用来管理容器之间关系的应用，并做到一个命令启动和关闭这些容器。



### 2.安装使用

#### 方法一

笔者下载有问题，可能是版本有问题。换了版本下载成功，但docker-compose 不存在

官网下载：https://docs.docker.com/compose/install/

```sh
curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

chmod +x /usr/local/bin/docker-compose

docker-compose --version
```

![image-20230319102216814](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319102216814.png)



#### 方法二

软件库进行安装（官方）

```sh
 $ sudo yum update
 $ sudo yum install docker-compose-plugin
```



### 3.三大核心概念和步骤

- 一文件：docker-compose.yml
- 两要素:
  - 服务:一个个应用容器实例，比如订单微服务、库存微服务、mysql容器、nginx容器或者redis容器
  - 工程:**由一组关联的应用容器组成的一个完整业务单元**，在 docker-compose.yml 文件中定义。



三个步骤

- 编写Dockerfile定义各个微服务应用并构建出对应的镜像文件
- 使用 docker-compose.yml 定义一个完整业务单元，安排好整体应用中的各个容器服务。
- 最后，执行**docker-compose up**命令 来启动并运行整个应用程序，完成一键部署上线；等价一次性运行多个docker run



```sh
Compose常用命令
docker-compose -h                           # 查看帮助
docker-compose up                           # 启动所有docker-compose服务
docker-compose up -d                        # 启动所有docker-compose服务并后台运行
docker-compose down                         # 停止并删除容器、网络、卷、镜像。
docker-compose exec  yml里面的服务id                 # 进入容器实例内部  docker-compose exec docker-compose.yml文件中写的服务id /bin/bash
docker-compose ps                      # 展示当前docker-compose编排过的运行的所有容器
docker-compose top                     # 展示当前docker-compose编排过的容器进程
 
docker-compose logs  yml里面的服务id     # 查看容器输出日志
docker-compose config     # 检查配置
docker-compose config -q  # 检查配置，有问题才有输出
docker-compose restart   # 重启服务
docker-compose start     # 启动服务
docker-compose stop      # 停止服务
 
```



#### 4.使用

前面讲到计算机网络的时候说到，容器宕机或者启停会导致ip地址发生变化，映射错误。

建议是通过服务名进行调用(不推荐写死IP)。

编写docker-compose.yml

```sh
version: "3"
 
services:
  microService:
    image: zzyy_docker:1.6
    container_name: ms01
    ports:
      - "6001:6001"
    volumes:
      - /app/microService:/data
    networks: 
      - atguigu_net 
     # 依赖的容器
    depends_on: 
      - redis
      - mysql
 
  redis:
    image: redis:6.0.8
    ports:
      - "6379:6379"
    volumes:
      - /app/redis/redis.conf:/etc/redis/redis.conf
      - /app/redis/data:/data
    networks: 
      - atguigu_net
     # 命令
    command: redis-server /etc/redis/redis.conf
 
  mysql:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: '123456'
      MYSQL_ALLOW_EMPTY_PASSWORD: 'no'
      MYSQL_DATABASE: 'db2021'
      MYSQL_USER: 'zzyy'
      MYSQL_PASSWORD: 'zzyy123'
    ports:
       - "3306:3306"
    volumes:
       - /app/mysql/db:/var/lib/mysql
       - /app/mysql/conf/my.cnf:/etc/my.cnf
       - /app/mysql/init:/docker-entrypoint-initdb.d
    networks:
      - atguigu_net
    command: --default-authentication-plugin=mysql_native_password #解决外部无法访问
 
networks: 
   atguigu_net: 
 

```



如此修改springboot的yml文件

```yaml
# ========================alibaba.druid相关配置=====================
spring.datasource.type=com.alibaba.druid.pool.DruidDataSource
spring.datasource.driver-class-name=com.mysql.jdbc.Driver
#spring.datasource.url=jdbc:mysql://192.168.111.169:3306/db2021?useUnicode=true&characterEncoding=utf-8&useSSL=false
# IP地址就换成服务名
spring.datasource.url=jdbc:mysql://mysql:3306/db2021?useUnicode=true&characterEncoding=utf-8&useSSL=false
spring.datasource.username=root
spring.datasource.password=123456
spring.datasource.druid.test-while-idle=false

# ========================redis相关配置=====================
spring.redis.database=0
#spring.redis.host=192.168.111.169
#host换成服务名
spring.redis.host=redis
spring.redis.port=6379
spring.redis.password=
spring.redis.lettuce.pool.max-active=8
spring.redis.lettuce.pool.max-wait=-1ms
spring.redis.lettuce.pool.max-idle=8
```



## 6.轻量级可视化工具Portainer

---

### 1.概念

> Portainer 是一款轻量级的应用，它提供了图形化界面，用于方便地管理Docker环境，包括单机环境和集群环境。



### 2.安装使用

**官网**

https://www.portainer.io/

https://docs.portainer.io/v/ce-2.9/start/install/server/docker/linux

**1. 命令行安装**

```sh
docker run -d -p 8000:8000 -p 9000:9000 --name portainer     --restart=always     -v /var/run/docker.sock:/var/run/docker.sock     -v portainer_data:/data     portainer/portainer

# --restart=always 表示该容器实例会随着docker的重启而重启
```

2. 访问地址 xxx:9090,创建admin用户

3. 登录后选择本地docker

![image-20230319194915131](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319194915131.png)



![image-20230319194941132](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319194941132.png)

对应的命令 `docker system df`

```sh
$ docker system df
TYPE            TOTAL     ACTIVE    SIZE      RECLAIMABLE
Images          10        5         1.947GB   980.6MB (50%)
Containers      5         1         4.368MB   4.367MB (99%)
Local Volumes   5         2         481.7MB   439.2MB (91%)
Build Cache     9         0         0B        0B
```



## 7.容器监控CAdvisor+InfluxDB+Granfana（重量级）

---

### 1.概念

对容器实例进行一个监控。原生的命令 `docker stats`

![image-20230319195855572](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319195855572.png)



缺点：

- 数据是实时展示，并没有地方进行存储和健康指标过线预警等功能



> CIG(CAdvisor监控收集+InfluxDB存储数据+Granfana展示图表)

![image-20230319200041128](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319200041128.png)

#### CAdvisor

![image-20230319200121895](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319200121895.png)

#### InfluxDB

![image-20230319200152037](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319200152037.png)

#### Granfana

![image-20230319200217278](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319200217278.png)



### 2.使用

通过docke compose 进行集中处理。

```yaml
version: '3.1'
 
volumes:
  grafana_data: {}
 
services:
 influxdb:
  image: tutum/influxdb:0.9
  restart: always
  environment:
  	# 提前创建一个数据库
    - PRE_CREATE_DB=cadvisor
  ports:
    - "8083:8083" #暴露端口
    - "8086:8086" #内部端口
  volumes:
    - ./data/influxdb:/data
 
 cadvisor:
  image: google/cadvisor
  links:
    - influxdb:influxsrv
  command: -storage_driver=influxdb -storage_driver_db=cadvisor -storage_driver_host=influxsrv:8086
  restart: always
  ports:
    - "8080:8080"
  volumes:
    - /:/rootfs:ro
    - /var/run:/var/run:rw
    - /sys:/sys:ro
    - /var/lib/docker/:/var/lib/docker:ro
 
 grafana:
  user: "104"
  image: grafana/grafana
  user: "104"
  restart: always
  links:
    - influxdb:influxsrv
  ports:
    - "3000:3000"
  volumes:
    - grafana_data:/var/lib/grafana
  environment:
    - HTTP_USER=admin
    - HTTP_PASS=admin
    - INFLUXDB_HOST=influxsrv
    - INFLUXDB_PORT=8086
    - INFLUXDB_NAME=cadvisor
    - INFLUXDB_USER=root
    - INFLUXDB_PASS=root

```



### 3.测试

- 浏览cAdvisor收集服务，http://ip:8080/。也有基础图形展示功能，主要用来数据采集。
- 浏览influxdb存储服务，http://ip:8083/
- 浏览grafana展现服务，http://ip:3000，默认帐户密码（admin/admin）

#### 配置数据源

![image-20230319202549114](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319202549114.png)

​	

![image-20230319202702471](C:/Users/wuxie/AppData/Roaming/Typora/typora-user-images/image-20230319202702471.png)



![image-20230319202710568](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319202710568.png)

账号密码都是`root`



#### 配置面板

![image-20230319202929844](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319202929844.png)

![image-20230319202936618](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319202936618.png)



![image-20230319202942997](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319202942997.png)



![image-20230319203033843](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319203033843.png)



![image-20230319203224942](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230319203224942.png)

