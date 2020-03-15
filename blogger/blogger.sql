/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.01
 Source Server Type    : MySQL
 Source Server Version : 50723
 Source Host           : localhost:3306
 Source Schema         : blogger

 Target Server Type    : MySQL
 Target Server Version : 50723
 File Encoding         : 65001

 Date: 15/03/2020 22:50:12
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '文章id',
  `category_id` bigint(20) unsigned NOT NULL COMMENT '分类id',
  `content` longtext NOT NULL COMMENT '文章内容',
  `title` varchar(1024) NOT NULL COMMENT '文章标题',
  `view_count` int(255) unsigned NOT NULL DEFAULT '0' COMMENT '阅读次数',
  `comment_count` int(255) unsigned NOT NULL DEFAULT '0' COMMENT '评论次数',
  `username` varchar(128) NOT NULL COMMENT '作者',
  `status` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '状态',
  `summary` varchar(256) NOT NULL COMMENT '文章摘要',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_view_count` (`view_count`) USING BTREE COMMENT '阅读次数索引',
  KEY `idx_comment_count` (`comment_count`) USING BTREE COMMENT '评论数索引',
  KEY `idx_category_id` (`category_id`) USING BTREE COMMENT '分类id索引'
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of article
-- ----------------------------
BEGIN;
INSERT INTO `article` VALUES (1, 1, 'hello,summer', 'summer', 10, 10, 'summer', 1, 'summer', '2020-03-13 17:47:04', '2020-03-13 17:47:04');
INSERT INTO `article` VALUES (2, 1, '2', '2', 2, 2, '2', 1, '2', '2020-03-13 18:47:18', '2020-03-13 18:47:18');
INSERT INTO `article` VALUES (3, 6, '<p>一般来说技术团队的金字塔顶尖往往是技术最牛的人做架构师（或TL）。所以架构师在广大码农中的占比大概平均不到 20%。而架构师也可以分为初级、中级、高级，江湖上真正高水平的架构师就更少了。</p><p>所以，80%码农干上许多年，还是做不了架构师，正在辛苦工作的程序员们，你有没有下面几种感觉？</p><p>&nbsp;</p><p>① 我的工作就是按时完成领导交给我的任务，至于代码写的怎样，知道有改进空间，但没时间去改进，关键是领导也不给时间啊。</p><p>&nbsp;</p><p>② 我发现我的水平总是跟不上技术的进步，有太多想学的东西要学，Swoft用的人最近比较多啊，听说最近Swoole比较火，还有微服务，听说PHP又更新了……</p><p>&nbsp;</p><p>③ 我发现虽然我工作5年了，除了不停的Coding写业务代码，Ctrl+c和Ctrl+V更熟练了，但编码水平并没有提高，还是一个普通程序员，但有人已经做到架构师了。</p><p>&nbsp;</p><p>④工作好几年了，想跳槽换个高薪工作，结果面试的考官都问了一些什么数据结构，什么垃圾回收，什么并发架构、协程编程之类的东西，虽然看过，但是平时用不着，看了也忘记了，回答不上来，结果面试官说我基础太差……</p><p>如果有以上问题，那么你绝对进入学习误区走了弯路；如果我们要成为架构师，我们自己要面临的三大问题：</p><p>找准定位：我是谁、我在哪里？<br>怎样做好架构师：我要做什么？<br>如何搭建架构师知识体系：我该怎么做？<br>如果你想要往架构师的方向发展的话，那或许你可以看一下我分享给你的这份进阶路线图，主要针对1到5年及以内的PHP开发人员，里面的技术包涵了PHP高并发、分布式、Swoole协程编程、微服务、框架内核、高性能等技术，这些也是目前中大型互联网企业比较常用的技术，那么来详细看看。</p><p>&nbsp;</p><p>一：常见模式与框架<br>学习PHP技术体系，设计模式，流行的框架</p><p>常见的设计模式，编码必备<br>Laravel、ThinkPHP开发必不可少的最新框架<br>YII、Symfony4.1核心源码剖析</p><p><br>二：微服务架构与性能优化<br>业务体系越来越复杂，协程编程，PHP并发编程、MySQL底层优化是架构升级的必经之路，PHP性能优化和微服务相关的技术有哪些呢？</p><p>Tars分布式RPC框架<br>Swoft微服务框架<br>服务器性能优化<br>算法与数据结构</p><p><br>&nbsp;</p><p>三：工程化与分布式架构<br>任何脱离细节的PPT架构师都是耍流氓，向上能运筹帷幄，向下能解决一线开发问题，PHP架构师需深入工程化、高并发，高可用，海量数据，没有分布式的架构知识肯定是玩不转的：</p><p>Linux操作/shell脚本编程<br>docker容器/自动化部署<br>分布式缓存/消息中间件<br>分布式架构原理/高并发分流</p><p><br>能掌握以上技术这些人必然具备在技术上独当一面的能力并且清楚自己未来的发展方向，从一个Coder逐步走向CTO或是架构师，成为项目组中不可或缺的人物。那么以上专题内容该怎么学习？为了大家能够顺利进阶中高级、架构师，我特地为大家准备了一套精品PHP架构师教程，保证你学了以后保证薪资上升一个台阶；</p><p>&nbsp;</p><p>对PHP后端技术，对PHP架构技术感兴趣的朋友，点击链接加入群聊【PHP高级交流群30】：621260710，一起学习，相互讨论。</p><p>群内已经有管理将知识体系整理好（源码，学习视频等资料），欢迎加群免费领取。</p><p>这套精品PHP教程绝不是市场上的那些妖艳贱货可比，作为web开发的佼佼者PHP并不逊色其他语言，加上Swoole后更加是如虎添翼！进军通信 、物联网行业开发百度地图、百度订单中心、虎牙、战旗TV等！寒冬裁员期过后正是各大企业扩大招人的时期，现在市场初级程序员泛滥，进阶中高级程序员绝对是各大企业急需的人才，这套学习教程适合那些1-5年以内的PHP开发者正处于瓶颈期，想要突破自己进阶中高级、架构师！名额有限，先到先得！</p><p>部分资料截图：</p><p>&nbsp;</p><p>&nbsp;</p><p>&nbsp;</p><p><br>还有限时精品福利：<br>★腾讯高级PHP工程师笔试题目 &nbsp; &nbsp;</p><p>★亿级PV高并发场景订单的处理</p><p>★laravel开发天猫商城组件服务</p><p>★战旗TV视频直播的架构项目实战</p><p>扫描下面二维码领取</p><p><br>对PHP后端技术，对PHP架构技术感兴趣的朋友，欢迎加QQ群：621260710，一起学习，相互讨论。</p><p>群内已经有管理将知识体系整理好（源码，学习视频等资料），欢迎加群免费领取。</p><p>本课程深度对标腾讯T3-T4标准，贴身打造学习计划为web开发人员进阶中高级、架构师提升技术，为自己增值涨薪！加入BAT特训营还可以获得内推大厂名额以及GO语言学习权限！！！</p>', '为什么80%的码农都做不了架构师？', 0, 0, '', 1, '<p>一般来说技术团队的金字塔顶尖往往是技术最牛的人做架构师（或TL）。所以架构师在广大码农中的占比大概平均不到 20%。而架构师也可以分为初级、中级、高级，江湖上真正高水平的架构师就更少了。</p><p>所以，80%码农干上许多年，还是做不了架构师，正在', '2020-03-13 23:41:19', '2020-03-13 23:41:19');
INSERT INTO `article` VALUES (4, 6, '<p>强调一下是<strong>我个人</strong>的见解以及接口在 <strong>Go 语言</strong>中的意义。</p><p>如果您写代码已经有了一段时间，我可能不需要过多解释接口所带来的好处，但是在深入探讨 Go 语言中的接口前，我想花一两分钟先来简单介绍一下接口。 如果您对接口很熟悉，请先跳过下面这段。</p><h2>接口的简单介绍</h2><p>在任一编程语言中，接口——方法或行为的集合，在功能和该功能的使用者之间构建了一层薄薄的抽象层。在使用接口时，并不需要了解底层函数是如何实现的，因为接口隔离了各个部分（划重点）。</p><p>跟不使用接口相比，使用接口的最大好处就是可以使代码变得简洁。例如，您可以创建多个组件，通过接口让它们以统一的方式交互，尽管这些组件的底层实现差异很大。这样就可以在编译甚至运行的时候动态替换这些组件。</p><p>用 Go 的 io.Reader 接口举个例子。io.Reader 接口的所有实现都有 Read(p []byte) (n int, err error) 函数。使用 io.Reader 接口的使用者不需要知道使用这个 Read 函数的时候那些字节从何而来。</p><h2>具体到 Go 语言</h2><p>在我使用 Go 语言的过程中，与我使用过的其他任何编程语言相比，我经常发现其他的、不那么明显的使用接口的原因。今天，我将介绍一个很普遍的，也是我遇到了很多次的使用接口的原因。</p><h2>Go 语言没有构造函数</h2><p>很多编程语言都有构造函数。构造函数是定义自定义类型（即 OO 语言中的类）时使用的一种建立对象的方法，它可以确保必须执行的任何初始化逻辑均已执行。</p><p>例如，假设所有 widgets 都必须有一个不变的，系统分配的标识符。在 Java 中，这很容易实现：</p><p>package io.krancour.widget;\r\n\r\nimport java.util.UUID;\r\n\r\npublic class Widget {\r\n\r\n &nbsp; &nbsp;private String id;\r\n\r\n &nbsp; &nbsp;// 使用构造函数初始化\r\n &nbsp; &nbsp;public Widget() {\r\n &nbsp; &nbsp; &nbsp; &nbsp;id = UUID.randomUUID().toString();\r\n &nbsp; &nbsp;}\r\n\r\n &nbsp; &nbsp;public String getId() {\r\n &nbsp; &nbsp; &nbsp; &nbsp;return id;\r\n &nbsp; &nbsp;}\r\n}\r\nclass App {\r\n &nbsp; &nbsp;public static void main( String[] args ){\r\n &nbsp; &nbsp; &nbsp; &nbsp;Widget w = new Widget();\r\n &nbsp; &nbsp; &nbsp; &nbsp;System.out.println(w.getId());\r\n &nbsp; &nbsp;}\r\n}\r\n</p><p>从上面这个例子可以看到，没有执行初始化逻辑就无法实例化一个新的 Widget 。</p><p>但是 Go 语言没有此功能。</p><figure class=\"image\"><img src=\"https://cdnjs.cloudflare.com/ajax/libs/emojify.js/1.1.0/images/basic/frowning.png\" alt=\":frowning:\"></figure><p>在 Go 语言中，可以直接实例化一个自定义类型。</p><p>定义一个 Widget 类型：</p><p>package widgets\r\n\r\ntype Widget struct {\r\n &nbsp; &nbsp;id string\r\n}\r\n\r\nfunc (w Widget) ID() string {\r\n &nbsp; &nbsp;return w.id\r\n}\r\n</p><p>可以像这样实例化和使用一个 widget：</p><p>package main\r\n\r\nimport (\r\n &nbsp; &nbsp;\"fmt\"\r\n &nbsp; &nbsp;\"github.com/krancour/widgets\"\r\n)\r\n\r\nfunc main() {\r\n &nbsp; &nbsp;w := widgets.Widget{}\r\n &nbsp; &nbsp;fmt.Println(w.ID())\r\n}\r\n</p><p>如果运行此示例，那么（也许）意料之中的结果是，打印出的 ID 是空字符串，因为它从未被初始化，而空字符串是字符串的“零值”。 我们可以在 widgets 包中添加一个类似于构造函数的函数来处理初始化：</p><p>package widgets\r\n\r\nimport uuid \"github.com/satori/go.uuid\"\r\n\r\ntype Widget struct {\r\n &nbsp; &nbsp;id string\r\n}\r\n\r\nfunc NewWidget() Widget {\r\n &nbsp; &nbsp;return Widget{\r\n &nbsp; &nbsp; &nbsp; &nbsp;id: uuid.NewV4().String(),\r\n &nbsp; &nbsp;}\r\n}\r\n\r\nfunc (w Widget) ID() string {\r\n &nbsp; &nbsp;return w.id\r\n}\r\n</p><p>然后我们简单地修改 main 来使用这个类似于构造函数的新函数：</p><p>package main\r\n\r\nimport (\r\n &nbsp; &nbsp;\"fmt\"\r\n &nbsp; &nbsp;\"github.com/krancour/widgets\"\r\n)\r\n\r\nfunc main() {\r\n &nbsp; &nbsp;w := widgets.NewWidget()\r\n &nbsp; &nbsp;fmt.Println(w.ID())\r\n}\r\n</p><p>执行该程序，我们得到了想要的结果。</p><p>但是仍然存在一个严重问题！我们的 widgets 包没有强制用户在初始一个 widget 的时候使用我们的构造函数。</p><h2>变量私有化</h2><p>首先我们尝试把自定义类型的变量私有化，以此来强制用户使用我们规定的构造函数来初始化 widget。在 Go 语言中，类型名、函数名的首字母是否大写决定它们是否可被其他包访问。名称首字母大写的可被访问（也就是 public ），而名称首字母小写的不可被访问（也就是 private ）。所以我们把类型 Widget 改为类型 widget ：</p><p>package widgets\r\n\r\nimport uuid \"github.com/satori/go.uuid\"\r\n\r\ntype widget struct {\r\n &nbsp; &nbsp;id string\r\n}\r\n\r\nfunc NewWidget() widget {\r\n &nbsp; &nbsp;return widget{\r\n &nbsp; &nbsp; &nbsp; &nbsp;id: uuid.NewV4().String(),\r\n &nbsp; &nbsp;}\r\n}\r\n\r\nfunc (w widget) ID() string {\r\n &nbsp; &nbsp;return w.id\r\n}\r\n</p><p>我们的 main 代码保持不变，这次我们得到了一个 ID 。这比我们想要的要近了一步，但是我们在此过程中犯了一个不太明显的错误。类似于构造函数的 NewWidget 函数返回了一个私有的实例。尽管编译器对此不会报错，但这是一种不好的做法，下面是原因解释。</p><p>在 Go 语言中，<i><strong>包</strong></i>是复用的基本单位。其他语言中的<i><strong>类</strong></i>是复用的基本单位。如前所述，任何无法被外部访问的内容实质上都是“包私有”，是该包的内部实现细节，对于使用这个包的使用者来说不重要。因此，Go 的文档生成工具 godoc 不会为私有的函数、类型等生成文档。</p><p>当一个公开的构造函数返回一个私有的 widget 实例，实际上就陷入了一条死胡同。调用这个函数的人哪怕有这个实例，也绝对在文档里找不到任何关于这个实例类型的描述，也更不知道 ID() 这个函数。Go 社区非常重视文档，所以这样做是不会被接受的。</p><h2>轮到接口上场了</h2><p>回顾一下，到目前为止，我们写了一个类似于构造函数的函数来解决 Go 语言缺乏构造函数的问题，但是为了确保人们用该函数而不是直接实例化 Widget ，我们更改了该类型的可见性——将其重命名为 widget，即私有化了。虽然编译器不会报错，但是文档中不会出现对这个私有类型的描述。不过，我们距离想要的目标还近了一步。接下来就要使用接口来完成后续的了。</p><p>通过创建一个<i><strong>可被访问的</strong></i>、widget 类型可以实现的接口，我们的构造函数可以返回一个公开的类型实例，并且会显示在 godoc 文档中。同时，这个接口的底层实现依然是私有的，使用者无法直接创建一个实例。</p><p>package widgets\r\n\r\nimport uuid \"github.com/satori/go.uuid\"\r\n\r\n// Widget is a ...\r\ntype Widget interface {\r\n &nbsp; &nbsp;// ID 返回这个 widget 的唯一标识符\r\n &nbsp; &nbsp;ID() string\r\n}\r\n\r\ntype widget struct {\r\n &nbsp; &nbsp;id string\r\n}\r\n\r\n// NewWidget() 返回一个新的 Widget 实例\r\nfunc NewWidget() Widget {\r\n &nbsp; &nbsp;return widget{\r\n &nbsp; &nbsp; &nbsp; &nbsp;id: uuid.NewV4().String(),\r\n &nbsp; &nbsp;}\r\n}\r\n\r\nfunc (w widget) ID() string {\r\n &nbsp; &nbsp;return w.id\r\n}\r\n</p><h2>总结</h2><p>我希望我已经充分地阐述了 Go 语言的这一特质——构造函数的缺失反而促进了接口的使用。</p><p>在我的下一篇文章中，我将介绍一种几乎与之相反的场景——在其他语言中要使用接口但是在 Go 语言中却不必。</p>', '在 Go 语言中，我为什么使用接口', 0, 0, '', 1, '<p>强调一下是<strong>我个人</strong>的见解以及接口在 <strong>Go 语言</strong>中的意义。</p><p>如果您写代码已经有了一段时间，我可能不需要过多解释接口所带来的好处，但是在深入探讨 Go 语言中的接口前，我想花一', '2020-03-13 23:44:30', '2020-03-13 23:44:30');
COMMIT;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `category_name` varchar(255) NOT NULL COMMENT '分类名字',
  `category_no` int(10) unsigned NOT NULL COMMENT '分类排序',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of category
-- ----------------------------
BEGIN;
INSERT INTO `category` VALUES (1, 'css/html', 1, '2018-08-12 10:55:45', '2018-08-12 10:59:00');
INSERT INTO `category` VALUES (2, '后端开发', 2, '2018-08-12 10:56:07', '2018-08-12 10:59:03');
INSERT INTO `category` VALUES (3, 'Java开发', 3, '2018-08-12 10:56:16', '2018-08-12 10:59:05');
INSERT INTO `category` VALUES (4, 'C++开发', 4, '2018-08-12 10:56:24', '2018-08-12 10:59:08');
INSERT INTO `category` VALUES (5, '架构剖析', 5, '2018-08-12 10:56:36', '2018-08-12 10:59:10');
INSERT INTO `category` VALUES (6, 'Golang开发', 6, '2018-08-12 10:56:45', '2018-08-12 10:59:14');
COMMIT;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '评论id',
  `content` text NOT NULL COMMENT '评论内容',
  `username` varchar(64) NOT NULL COMMENT '评论作者',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '评论发布时间',
  `status` int(255) unsigned NOT NULL DEFAULT '0' COMMENT '评论状态: 0, 删除；1， 正常',
  `article_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of comment
-- ----------------------------
BEGIN;
INSERT INTO `comment` VALUES (1, '', 'summer', '2020-03-14 21:59:31', 0, 0);
INSERT INTO `comment` VALUES (2, '', 'summer', '2020-03-14 22:00:15', 0, 0);
INSERT INTO `comment` VALUES (3, '', 'summer', '2020-03-14 22:00:46', 0, 0);
INSERT INTO `comment` VALUES (4, '', 'summer', '2020-03-14 22:01:16', 0, 0);
INSERT INTO `comment` VALUES (5, '', 'summer', '2020-03-14 22:06:39', 0, 4);
INSERT INTO `comment` VALUES (6, 'hello,world', 'summer', '2020-03-14 22:08:12', 0, 4);
COMMIT;

-- ----------------------------
-- Table structure for leave
-- ----------------------------
DROP TABLE IF EXISTS `leave`;
CREATE TABLE `leave` (
  `id` bigint(20) NOT NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
