# Spring深度学习

# 1.原生servlet

---

## 1.准备工作

用 **Servlet3.0+** 的版本可以基于注解开发，效率较快）

```xml
<dependencies>
    <dependency>
        <groupId>javax.servlet</groupId>
        <artifactId>javax.servlet-api</artifactId>
        <version>3.1.0</version>
        <!-- 这个scope 只能作用在编译和测试时,同时没有传递性--->
        <scope>provided</scope>
    </dependency>
</dependencies>

```

引入maven编译插件

```xml
<build>
<plugins>
    <plugin>
<!--- maven编译插件 --->           <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-compiler-plugin</artifactId>
        <version>3.2</version>
        <configuration>
            <source>1.8</source>
            <target>1.8</target>
            <encoding>UTF-8</encoding>
        </configuration>
    </plugin>
</plugins>
</build>
```

打包方式为`war`

 ```xml
 <packaging>war</packaging>
 ```



配置部署环境tomat

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/31aa147ff1cb4f56be12b111861c7306~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/14ac30da10b84813be0abbadbe6f83b4~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



## 2.编写测试

```java
@WebServlet(urlPatterns = "/demo1")
public class DemoServlet1 extends HttpServlet {
    
    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {
        response.getWriter().println("DemoServlet1 run ......");
    }
    
}
```

输入[localhost:8080/Spring_Depth_Learning_war_exploded/demo1](http://localhost:8080/Spring_Depth_Learning_war_exploded/demo1)

![image-20230321095047214](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230321095047214.png)



注意：如果控制台出现乱码，修改`tomcat/conf/logging.properties` 

```properties
java.util.logging.ConsoleHandler.encoding = GBK(原UTF-8)
```



## 3.编写MVC框架

![img](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/bd19682e7682438fb4ddd1460342ac23~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



**代码框架如下**

![image-20230321095602151](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/image-20230321095602151.png)

1. Dao与DaoImpl

```java
public interface DemoDao {
    List<String> findAll();
}

===============
public class DemoDaoImpl implements DemoDao {
    
    @Override
    public List<String> findAll() {
        // 此处应该是访问数据库的操作，用临时数据代替
        return Arrays.asList("mysql", "mysql", "mysql");
    }
}
```

2. Service与ServiceImpl

```java
public interface DemoService {
    List<String> findAll();
}

=============
public class DemoServiceImpl implements DemoService {
    
    private DemoDao demoDao = new DemoDaoImpl();
    
    @Override
    public List<String> findAll() {
        return demoDao.findAll();
    }
}
```

3. Servlet

```java
@WebServlet(urlPatterns = "/demo1")
public class DemoServlet1 extends HttpServlet {
    
    DemoService demoService = new DemoServiceImpl();
    
    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
        resp.getWriter().println(demoService.findAll().toString());
    }
}
```



## 4.问题和解决

### 问题1：

> 修改数据库资源，从MySQL切换到Oracle。

```java
public class DemoDaoImpl implements DemoDao {
    
    @Override
    public List<String> findAll() {
        // 此处应该是访问数据库的操作，用临时数据代替
        //return Arrays.asList("mysql", "mysql", "mysql");
        return Arrays.asList("oracle", "oracle", "oracle")
    }
}
```

但如果我们又要使用mysql就需要再次改动源码，我们应该避免这种操作的发生。

**解决：**

> **使用静态工厂**，将不同数据源的dao层写好，然后通过静态工厂调用，就可以避免源码的改动,只需要修改调用即可。

```java
public class BeanFactory {
    public static DemoDao getDemoDao() {
        // return new DemoDaoImpl();
        return new DemoOracleDao();
    }
}
```

改动ServiceImpl

```java
public class DemoServiceImpl implements DemoService {
    
    DemoDao demoDao = BeanFactory.getDemoDao();
    
    @Override
    public List<String> findAll() {
        return demoDao.findAll();
    }
}
```



### 问题2：

> 类与类之间耦合度较高

```java
public class BeanFactory {
    public static DemoDao getDemoDao() {
        // return new DemoDaoImpl();
        return new DemoOracleDao();
    }
}
```

如上可以描述为`BeanFactory`的`getDemoDao` 使用**强依赖于**`DemoOracleDao`。如果`DemoOracleDao` 不存在，那么就导致整个项目编译失败，也称之为**"紧耦合"**

解决：

> 反射

```java
public class BeanFactory {
    
    public static DemoDao getDemoDao() {
        try {
            return (DemoDao) Class.forName("com.linkedbear.architecture.c_reflect.dao.impl.DemoDaoImpl").newInstance();
        } catch (Exception e) {
            e.printStackTrace();
            throw new RuntimeException("DemoDao instantiation error, cause: " + e.getMessage());
        }
    }
}
```

因为反射找到类不存在，就会抛出`ClassNotFoundException` 异常，但项目是正常启动了的。两个类之间的耦合度就降低，称为**"弱依赖"**



### 问题3：

> 硬编码

问题2代码可见，在使用反射时，我们的全类路径是写死了的，如果要该数据源dao层就要重新输入全类路径

解决

> 使用外部配置文件

在`src/main/resource` 下创建`factory.properties`

```properties
demoService=com.linkedbear.architecture.d_properties.service.impl.DemoServiceImpl
demoDao=com.linkedbear.architecture.d_properties.dao.impl.DemoDaoImpl
```



```java
public class BeanFactory {
    
    private static Properties properties;
     // 使用静态代码块初始化properties，加载factord.properties文件
    static {
        properties = new Properties();
        try {
            // 必须使用类加载器读取resource文件夹下的配置文件
            properties.load(BeanFactory.class.getClassLoader().getResourceAsStream("factory.properties"));
        } catch (IOException e) {
            // BeanFactory类的静态初始化都失败了，那后续也没有必要继续执行了
            throw new ExceptionInInitializerError("BeanFactory initialize error, cause: " + e.getMessage());
        }
    }
 ================================
     public static DemoDao getDemoDao() {
        try {
            Class<?> beanClazz = Class.forName(properties.getProperty("demoDao"));
            return beanClazz.newInstance();
        } catch (ClassNotFoundException e) {
            throw new RuntimeException("BeanFactory have not [" + beanName + "] bean!", e);
        } catch (IllegalAccessException | InstantiationException e) {
            throw new RuntimeException("[" + beanName + "] instantiation error!", e);
        }
    }
```

**并且这里的"demoDao"也不用写死，通过参数传递**

```java
  public static Object getBean(String beanName) {
        try {
            // 从properties文件中读取指定name对应类的全限定名，并反射实例化
            Class<?> beanClazz = Class.forName(properties.getProperty(beanName));
            return beanClazz.newInstance();
        } catch (ClassNotFoundException e) {
            throw new RuntimeException("BeanFactory have not [" + beanName + "] bean!", e);
        } catch (IllegalAccessException | InstantiationException e) {
            throw new RuntimeException("[" + beanName + "] instantiation error!", e);
        }
```



业务层使用

```java
public class DemoServiceImpl implements DemoService {
    
    DemoDao demoDao = (DemoDao) BeanFactory.getBean("demoDao");
```



### 问题4

> return beanClazz.newInstance(); 

导致每次使用时都会新创建对象，对此是没有必要和浪费资源的

解决

> 使用缓存



```java
public class BeanFactory {
    // 缓存区，保存已经创建好的对象
    private static Map<String, Object> beanMap = new HashMap<>();
    
    // ......
    
    public static Object getBean(String beanName) {
    // 双检锁保证beanMap中确实没有beanName对应的对象
    if (!beanMap.containsKey(beanName)) {
        synchronized (BeanFactory.class) {
            if (!beanMap.containsKey(beanName)) {
                // 过了双检锁，证明确实没有，可以执行反射创建
                try {
                    Class<?> beanClazz = Class.forName(properties.getProperty(beanName));
                    Object bean = beanClazz.newInstance();
                    // 反射创建后放入缓存再返回
                    beanMap.put(beanName, bean);
                } catch (ClassNotFoundException e) {
                    throw new RuntimeException("BeanFactory have not [" + beanName + "] bean!", e);
                } catch (IllegalAccessException | InstantiationException e) {
                    throw new RuntimeException("[" + beanName + "] instantiation error!", e);
                }
            }
        }
    }
    return beanMap.get(beanName);
}
```



#### 总结

- 静态工厂可将多处依赖抽取分离
- 外部化配置文件+反射可解决配置的硬编码问题
- 缓存可控制对象实例数

也是引入了spring的一个重要思想**"控制反转（ Inverse of Control , IOC )"**,将对象的创建交给工厂去做，并且根据`beanName`去获得和创建对象，这个过程称为**"依赖查找(Dependency Lookup , DL)"**



# 2.spring框架

## 介绍

> Spring 框架为**任何类型的部署平台**上的**基于 Java** 的现代**企业应用程序**提供了全面的**编程和配置模型**。
>
> Spring 的一个关键元素是在**应用程序级别的基础架构支持**：Spring 专注于企业应用程序的 “**脚手架**” ，以便团队可以**专注于应用程序级别的业务逻辑**，而不必与特定的部署环境建立不必要的联系。
>
> - 任何类型的部署平台：无论是操作系统，还是 Web 容器（ Tomcat 等）都是可以部署基于 SpringFramework 的应用
> - 企业应用程序：包含 JavaSE 和 JavaEE 在内，它被称为一站式解决方案
> - 编程和配置模型：基于框架编程，以及基于框架进行功能和组件的配置
> - 基础架构支持：SpringFramework 不含任何业务功能，它只是一个**底层的应用抽象支撑**
> - 脚手架：使用它可以更快速的构建应用



> SpringFramework 是一个**容器框架**，它集成了各个类型的工具，通过核心的 IOC 容器实现了底层的组件实例化和生命周期管理。
>
> - IOC & AOP：SpringFramework 的两大核心特性：**Inverse of Control 控制反转、Aspect Oriented Programming 面向切面编程**
> - 轻量级：对比于重量级框架，它的规模更小（可能只有几个 jar 包）、消耗的资源更少
> - 一站式：覆盖企业级开发中的所有领域
> - 第三方整合：SpringFramework 可以很方便的整合进其他的第三方技术（如持久层框架 MyBatis / Hibernate ，表现层框架 Struts2 ，权限校验框架 Shiro 等）
> - 容器：SpringFramework 的底层有一个管理对象和组件的容器，由它来支撑基于 SpringFramework 构建的应用的运行



## 相关面试题

### 1.spring框架是什么，如何进行概述

**SpringFramework 是一个开源的、松耦合的、分层的、可配置的一站式企业级 Java 开发框架，它的核心是 IOC 与 AOP ，它可以更容易的构建出企业级 Java 应用，并且它可以根据应用开发的组件需要，整合对应的技术。**



分层:

![img](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/02/cda2894078094fe5833bb8838ffca214.png)



## 2.为什么使用spring框架

- **IOC**：组件之间的解耦（咱上一章已经体会到了）
- **AOP**：切面编程可以将应用业务做统一或特定的功能增强，能实现应用业务与增强逻辑的解耦
- **容器**与事件：管理应用中使用的组件Bean、托管Bean的生命周期、事件与监听器的驱动机制
- Web、事务控制、测试、与**其他技术的整合**



## 3.spring包含的模块

- beans、core、context、expression 【核心包】
- aop 【切面编程】
- jdbc 【整合 jdbc 】
- orm 【整合 ORM 框架】
- tx 【事务控制】
- web 【 Web 层技术】
- test 【整合测试】
- ......



# 3.spring使用

## 1.IOC-DL（依赖查询)

### 1.IOC-DL 入门（byName）

#### 1.1引入依赖

~~~xml
<dependency>
    <groupId>org.springframework</groupId>
    <artifactId>spring-context</artifactId>
    <version>5.2.8.RELEASE</version>
</dependency>
~~~

#### 1.2创建类

在`org.example` 下创建一个bean包，随后创建一个`Person` 类

![image-20230410223919202](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/02/image-20230410223919202.png)

#### 1.3创建配置文件

在`resource` 下创建一个配置文件，IOC可以通过配置文件来扫描得到类，对象的信息。

这里`bean`的配置，`id`唯一标识来获取bean对象，`class` 对应创建bean的类

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">
    <bean id="person" class="org.example.AnnotationBean.Person"></bean>
</beans>
```

#### 1.4 启动类

```java
public class Main {
    public static void main(String[] args) {
        BeanFactory factory = new ClassPathXmlApplicationContext("quickstart-byname.xml");
        Person person = (Person) factory.getBean("person");
        System.out.println(person);
        //org.example.AnnotationBean.Person@4116aac9

    }
}
```

通过`ClassPathXmlApplicationContext` 来加载配置文件，通过`BeanFactory` 向上继承接口收，然后`getBean (id 配置文件中的)` ，获得bean对象，接着就可以输出了。

运行 `main` 方法，可以成功打印出 `Person` 的全限定类名 + 内存地址，证明编写成功



### 2.IOC-DL(byType)

上面依赖查找通过id进行获取，还可以通过**类型查找**

#### 2.1不声明id属性

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">
    <bean  class="org.example.AnnotationBean.Person"></bean>
</beans>
```

#### 2.2 启动

```java
public class Main {
    public static void main(String[] args) {
        BeanFactory factory = new ClassPathXmlApplicationContext("quickstart-byname.xml");
        Person person =  factory.getBean(Person.class);
        System.out.println(person);
        //org.example.AnnotationBean.Person@489115ef
    }
}
```



### 3.IOC-DL（接口与实现类）

IOC还可以根据接口来查找实现类。将前面的dao层拿来

![image-20230410225108019](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/02/image-20230410225108019.png)

#### 3.1 新增bean

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">
    <bean  class="org.example.AnnotationBean.Person"></bean>
    <bean class="org.example.dao.impl.DemoDaoImpl"></bean>
</beans>
```



#### 3.2 启动类

```java
public class Main {
    public static void main(String[] args) {
        BeanFactory factory = new ClassPathXmlApplicationContext("quickstart-byname.xml");
        DemoDao bean = factory.getBean(DemoDao.class);
        System.out.println(bean.findAll());
        // [aaa, bbb, ccc]
    }
}
```

**由上可得通过接口来获得实现类**



### 4.IOC-DL (getBeansOfType)

问题抛出：如果一个接口有多个实现类，想要一次性拿出，那么`getBean` 就不够用，可以使用`getBeansOfType` 的查找方式。

#### 4.1 声明Bean和配置文件

![image-20230410232651750](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/02/image-20230410232651750.png)



```xml
<bean class="org.example.dao.impl.DemoDaoImpl"/>
<bean class="org.example.dao.impl.DemoMySQLImpl"/>
<bean class="org.example.dao.impl.DemoOracleImpl"/>
```



#### 4.2 测试启动类

![image-20230410233003747](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/02/image-20230410233003747.png)

可见没有`getBeansOfType` 这个方法。固然是接口使用错误

将`BeanFactory` 切换为`ApplicationContext`



```java
Map<String, DemoDao> beansOfType = factory.getBeansOfType(DemoDao.class);
beansOfType.forEach((name, bean) -> {
    System.out.println(name + ":" + bean);
});

/*
org.example.dao.impl.DemoDaoImpl#0:org.example.dao.impl.DemoDaoImpl@59af0466
org.example.dao.impl.DemoMySQLImpl#0:org.example.dao.impl.DemoMySQLImpl@3e6ef8ad
org.example.dao.impl.DemoOracleImpl#0:org.example.dao.impl.DemoOracleImpl@346d61be
*/
```

**这样就实现了返回所有实现类**



#### 4.3  BeanFactory与ApplicationContext

![image-20230410233559912](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/02/image-20230410233559912.png)

可以看到`ApplicationContext` 是`BeanFactory` 的子接口。



##### 4.3.1 官方解释

> `org.springframework.beans` 和 `org.springframework.context` 包是 SpringFramework 的 IOC 容器的基础。`BeanFactory` 接口提供了一种高级配置机制，能够管理任何类型的对象。`ApplicationContext` 是 `BeanFactory` 的子接口。它增加了：
>
> - 与 SpringFramework 的 AOP 功能轻松集成
> - 消息资源处理（用于国际化）
> - 事件发布
> - 应用层特定的上下文，例如 Web 应用程序中使用的 `WebApplicationContext`
>
> 你应该使用 `ApplicationContext` ，除非能有充分的理由解释不需要的原因。一般情况下，我们推荐将 `GenericApplicationContext` 及其子类 `AnnotationConfigApplicationContext` 作为自定义引导的常见实现。这些实现类是用于所有常见目的的 SpringFramework 核心容器的主要入口点：加载配置文件，触发类路径扫描，编程式注册 Bean 定义和带注解的类，以及（从5.0版本开始）注册功能性 Bean 的定义。

总之一句话`ApplicationContext` 真的比 `BeanFactory` 强大太多了，所以咱还是选择使用 `ApplicationContext` 吧！



##### 4.3.2 面试题

##### BeanFactory与ApplicationContext的对比

> `BeanFactory` 接口提供了一个**抽象的配置和对象的管理机制**，`ApplicationContext` 是 `BeanFactory` 的子接口，它简化了与 AOP 的整合、消息机制、事件机制，以及对 Web 环境的扩展（ `WebApplicationContext` 等），`BeanFactory` 是没有这些扩展的。



WebApplicationContext 可以加载包含 Web 应用特定的 bean 的 XML 文件，如**控制器、视图解析器和处理器映射器**等。此外，WebApplicationContext 还提供了许多有用的功能，比如国际化和主题支持等。在 Spring MVC 中，DispatcherServlet 就会创建一个 WebApplicationContext 来处理请求，并将其委托给相应的处理器来处理请求。

### 5.withAnnotation

**还可以根据标注的注解来查找对应的Bean** 

#### 5.1 声明bean+注解+配置文件

##### 5.1.1 新建一个`Annotaion` 包，声明一个注解:`@Color` 

```java
@Documented
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface Color {
}

```

解读注解

1.  @Documented 是一个标记注解，用于指示该注解应该包含在Javadoc中生成的文档中。当使用 javadoc 命令生成文档时，@Documented 修饰的注解会被包括在文档中
2.  @Target(ElementType.TYPE) 表示该**注解只能用于类、接口或枚举类型**上。如果将该注解用于其他类型（如方法或字段），编译器将报错
3. @Retention(RetentionPolicy.RUNTIME) 表示该**注解在运行时保留**，可以通过反射机制获取注解信息。



#### 5.1.2 **配置文件记得加上**

```xml
    <bean class="org.example.bean.Black"/>
    <bean class="org.example.bean.Red"/>
```



#### 5.1.3 测试启动类

使用`getBeansWithAnnotaion` 方法进行获取

```java
        Map<String, Object> beansWithAnnotation = factory.getBeansWithAnnotation(Color.class);
        beansWithAnnotation.forEach((name, bean) -> {
            System.out.println(name + ":" + bean);
        });
//org.example.bean.Black#0:org.example.bean.Black@6babf3bf
//org.example.bean.Red#0:org.example.bean.Red@3059cbc
```



### 6.获得IOC容器的所有Bean

`getBeanDefinitionNames` 获得所有Bean的id

测试启动类,可得打印了创建的bean对象。定义了id的打印id，没有定义则是输出对应全类名

```java
//5.获得全部bean
String[] beanDefinitionNames = factory.getBeanDefinitionNames();
Stream.of(beanDefinitionNames).forEach(System.out::println);

/*
person
org.example.AnnotationBean.Cat#0
org.example.dao.impl.DemoDaoImpl#0
org.example.dao.impl.DemoMySQLImpl#0
org.example.dao.impl.DemoOracleImpl#0
org.example.bean.Black#0
org.example.bean.Red#0
*/
```



### 7.延迟查询

对于一些特殊的场景，需要依赖容器中的某些特定的 Bean ，但当它们不存在时也能使用默认 / 缺省策略来处理逻辑。



#### 7.1 使用现有方案进行缺省加载

我们再创建一个`Green` 类，但是不在配置文件里面注册直接使用，调用`getBean` 时就会报错。

`Exception in thread "main" org.springframework.beans.factory.NoSuchBeanDefinitionException: No qualifying bean of type 'org.example.bean.Green' available`

此时可以使用异常处理`try...catch.` 进行手动创建

```java
Green bean = null;
try {
    bean = factory.getBean(Green.class);
} catch (BeansException e) {
    bean = new Green();
}
```

但这样并不优美，并且bean多时，代码量增加。



#### 7.2 优化-获取先进行检查

```java
Green green = factory.containsBean("green") ? (Green) factory.getBean("green") : new Green();
```

缺点：`containsBean` 只能传`bean`的id，不能查类型。所以不存在bean时，id不知道是什么，就导致方法问题。



#### 7.3 优化-延迟查找

[延迟注入](#2.8.2 延迟注入)

上面获取bean时，都是在编译阶段查找是否存在，如果不存在就报错了。那么如果有些bean没有使用，也会导致程序报错。所以我们想**是否能不直接报错，而是在我使用的时候才检查是否存在再报错**。**`ObjectProvider`**   便可以实现延迟查找。

```java
ObjectProvider<Green> beanProvider = factory.getBeanProvider(Green.class);
```

可以发现，`ApplicationContext` 中有一个方法叫 `getBeanProvider` ，它就是返回类似**“包装”**。如果直接 `getBean` ，那如果容器中没有对应的 Bean ，就会报 `NoSuchBeanDefinitionException`；如果使用这种方式，运行 `main` 方法后发现并没有报错，只有调用 `dogProvider` 的 `getObject` ，真正要取包装里面的 Bean 时，才会报异常。所以总结下来，`ObjectProvider` 相当于**延后了 Bean 的获取时机，也延后了异常可能出现的时机**。



但是还是没有解决问题，调用`getObject` 还会报错。



#### 7.4 延迟查找-方案实现

`ObjectProvider` 中还有一个方法：`getIfAvailable` ，它可以在**找不到 Bean 时返回 null 而不抛出异常**。使用这个方法，就可以避免上面的问题了。改良之后的代码如下：

```java
ObjectProvider<Green> beanProvider = factory.getBeanProvider(Green.class);
Green green = beanProvider.getIfAvailable();
if (green == null) {
    green = new Green();
}
```



#### 7.5 `ObjectProvider` 在jdk8的升级

`ObjectProvider` 在 SpringFramework 5.0 后扩展了一个带 `Supplier` 参数的 `getIfAvailable` ，它可以在找不到 Bean 时直接用 **`Supplier`** 接口的方法返回默认实现，由此上面的代码还可以进一步简化为：

```java
beanProvider.getIfAvailable(() -> new Green());
beanProvider.getIfAvailable(Green::new);
```



`ifAvailable` 可以在bean存在时执行相关方法

```java
public class Green {

    public void say() {
        System.out.println("hello");
    }
}
//========================
beanProvider.ifAvailable(Green::say); //hello(记得注册)
```





## 2.IOC-DI(依赖注入)

如上创建的bean都是不带属性，如果想要预设属性就需要使用到IOC的**依赖注入**。



### 2.1 简单demo

#### 2.1.1 声明类,编写注解

这里使用`Lombok`注解

```xml
<dependency>
    <groupId>org.projectlombok</groupId>
    <artifactId>lombok</artifactId>
    <version>1.18.26</version>
</dependency>
```



```java
@Data
public class Person {
	String name;
    int age;
}
```

启动类查看结果

```java
public class Main {
    public static void main(String[] args) {
        BeanFactory factory = new ClassPathXmlApplicationContext("quickstart-byname.xml");
        Person person =  factory.getBean(Person.class);
//        DemoDao bean = factory.getBean(DemoDao.class);
//        System.out.println(bean.findAll());
        
        System.out.println(person);
		//Person(name=null, age=0)
    }
}
```



#### 2.1.2 赋值

在`<bean>` 标签里面可以声明`property` 标签，来进行属性赋值

```xml
<bean class="org.example.AnnotationBean.Person">
    <property name="name" value="思无邪"></property>
    <property name="age" value="20"></property>
</bean>
```

1. `name` 对应字段名
2. `value` 对应字段属性

启动结果

```sh
Person(name=思无邪, age=20)
```



### 2.2 关联Bean赋值

#### 2.2.1 新增类

再创建一个类,引用另一个类`Person`,此时属性赋值就可以使用`ref` 属性，代表**关联赋值的Bean的id**

```java
@Data
public class Cat {
    Person master;
    String name;
}
```



#### 2.2.2 编写配置文件

```xml
<bean class="org.example.AnnotationBean.Person" id="person">
    <property name="name" value="思无邪"></property>
    <property name="age" value="20"></property>
</bean>

<bean class="org.example.AnnotationBean.Cat">
    <property name="name" value="小黄" ></property>
    <!-- ref = 对应上面Person Bean的id    -->
    <property name="master" ref="person"></property>
</bean>
```



#### 2.2.3 启动类查看结果

```sh
Person(name=思无邪, age=20)
Cat(master=Person(name=思无邪, age=20), name=小黄)
```



### 2.3 setter属性注入

#### 2.3.1 xml方式的setter注入

```xml
<bean id="person" class="com.linkedbear.spring.basic_di.a_quickstart_set.bean.Person">
    <property name="name" value="test-person-byset"/>
    <property name="age" value="18"/>
</bean>
```

#### 2.3.2 注解方式的setter注入

```java
@Bean
public Person person() {
    Person person = new Person();
    person.setName("test-person-anno-byset");
    person.setAge(18);
    return person;
}
```

### 2.4 构造器注入

一些bean的属性依赖，在调用构造器（构造方法）时就设置好；获得有些bean没有无参构造器就必须要使用**构造器注入**

#### 2.4.1 修改Bean

新增全参构造器

```java
public Person(String name, int age) {
    this.name = name;
    this.age = age;
}
//============ 或者Lombok的注解
@AllArgsConstructor
```

默认的无参构造器就无了，这样xml注册时`<bean>` 标签构建时就失效，提示没有默认的构造方法。

```sh
Caused by: java.lang.NoSuchMethodException: com.linkedbear.spring.basic_di.b_constructor.bean.Person.<init>()
```

#### 2.4.2 xml方式的构造器注入

```xml
<bean class="org.example.AnnotationBean.Person">
    <constructor-arg index="0" value="思无邪"/>
    <constructor-arg index="1" value="19"/>
</bean>
```



#### 2.4.3 注解式构造器属性注入

```java
@Bean
public Person person() {
    return new Person("思无邪", 18);
}
```



### 2.5 注解式属性注入

因为注册bean的方式不仅有`@Bean`,还有组件扫描。所以后者方式，就需要进行另外的方法注解式注入。



#### 2.5.1 @Component下的属性注入

简单的属性注入方式:`@Value`

```java
@Component
@Data
public class Green {
    @Value("hello")
    private String name;
    @Value("1")
    private Integer order;

    @Override
    public String toString() {
        return "Green{" +
                "name='" + name + '\'' +
                ", order=" + order +
                '}';
    }

}
```



启动类测试得

```java
public class AnnotationConfigApplication {
    public static void main(String[] args) {
        ApplicationContext ctx = new AnnotationConfigApplicationContext("org.example.bean");
        Green bean = ctx.getBean(Green.class);
        System.out.println(bean);
    }
}
//Green{name='hello', order=1}
```



#### 2.5.2 外部配置文件引入-@PropertySource

之前读取property配置文件，使用的是Properties类去IO读取。spring提供了`@PropertySource` 注解来导入外部的配置文件



##### 2.5.2.1 创建bean+配置文件

新建一个`Blue`类

```java
@Component
@Data
public class Blue {
    
    private String name;
    
    private Integer order;

    @Override
    public String toString() {
        return "Blue{" +
                "name='" + name + '\'' +
                ", order=" + order +
                '}';
    }
}

```

在resource下创建`Blue.properties`文件，用户存放Blue类的相关属性

```properties
blue.name=hello
blue.order=1
```

##### 2.5.2.2 属性注入

然后在类里面使用`@Vlue` 获得数据

```java
@Component
@Data
public class Blue {

    @Value("${blue.name}")
    private String name;
    @Value("${blue.order}")
    private Integer order;
}
```



启动类测试得

```sh
Blue{name='hello', order=1}
```



##### 2.5.2.3 xml中使用占位符

用法与`@Vlue` 一致

```xml
<!--- 注明配置文件 -->
<context:property-placeholder ignore-unresolvable="true" location="Blue.properties"/> 

<bean class="org.example.bean.Blue">
        <property name="name" value="${blue.name}"/>
        <property name="order" value="${blue.order"/>
</bean>
```

启动类测试得

```java
ApplicationContext ctx = new ClassPathXmlApplicationContext("application.xml");
Blue bean = ctx.getBean(Blue.class);
System.out.println(bean);
//Blue{name='hello', order=1}
```



说明：

`ignore-unresolvable`

> 当设置为 `true` 时，如果找不到与给定属性名匹配的属性值，则会忽略该属性，并且不会抛出任何异常。这对于需要配置大量属性的应用程序非常有用，因为它可以确保即使在一些属性没有被正确设置的情况下，应用程序也可以正常运行。
>
> 当设置为 `false` 时，如果找不到与给定属性名匹配的属性值，则会抛出 `IllegalArgumentException` 异常。这可能会导致应用程序无法启动或运行，因此需要谨慎使用。
>
> 在Spring Boot中，该属性默认设置为 `true`，因此在属性文件中使用不存在的属性名称不会引发异常。



#### 2.5.3 SpEL表达式

如一个 Bean 需要依赖另一个 Bean 的某个属性，或者需要动态处理一个特定的属性值，这种情况**${}** 占位符实现不了（占位符只能取配置项),就需要使用到**SpEL表达式**



##### 2.5.3.1 概念

> SpEL 全称 Spring Expression Language ，它从 SpringFramework 3.0 开始被支持，它本身可以算 SpringFramework 的组成部分，但又可以被独立使用。它可以支持调用属性值、属性参数以及方法调用、数组存储、逻辑计算等功能。



##### 2.5.3.2 属性注入

SpEL的语法统一用`#{}` 表示。

创建新bean

```java
@Component
@Data
public class Yellow {

    @Value("#{'黄色'}")
    private String name;

    @Value("#{1}")
    private Integer order;

    @Override
    public String toString() {
        return "Yellow{" +
                "name='" + name + '\'' +
                ", order=" + order +
                '}';
    }
}
```



启动类测试得

```java
public class AnnotationConfigApplication {
    public static void main(String[] args) {
        ApplicationContext ctx = new AnnotationConfigApplicationContext(ComponentScanConfiguration.class);
        Yellow bean = ctx.getBean(Yellow.class);
        System.out.println(bean);
        //Yellow{name='黄色', order=1}
    }
}
```





##### 2.5.3.3 Bean属性引用

```java
@Component
@Data
public class Yellow {

    @Value("#{'copy of ' + blue.name}")
    private String name;

    @Value("#{blue.order + 1}")
    private Integer order;

    @Override
    public String toString() {
        return "Yellow{" +
                "name='" + name + '\'' +
                ", order=" + order +
                '}';
    }
}

// Yellow{name='copy of hello', order=2}
```



xml使用方式:

```xml
<bean class="com.linkedbear.spring.basic_di.c_value_spel.bean.Green">
    <property name="name" value="#{'copy of ' + blue.name}"/>
    <property name="order" value="#{blue.order + 1}"/>
</bean>
```



##### 2.5.3.4 方法调用

新建一个类`White`

```java
@Component
@Data
public class White {

    @Value("#{yellow.name.substring(0,3)}")
    private String name;

    @Value("#{T(java.lang.Integer).MAX_VALUE}")
    private Integer order;

}
```



```java
public class AnnotationConfigApplication {
    public static void main(String[] args) {
        ApplicationContext ctx = new AnnotationConfigApplicationContext(ComponentScanConfiguration.class);
        White bean = ctx.getBean(White.class);
        System.out.println(bean);
        //White{name='cop', order=2147483647}
    }
}
```



### 2.6 自动注入

xml里面可以使用`ref` 在一个bean里面注入另一个bean，注解当然可以做到。

#### 2.6.1 @Autowired

#####  2.6.1.1创建类

```java
@Data
@Component
public class Person {

    private String name = "思无邪";
}

```

`Cat`

```java
@Data
public class Cat {

    Person master;
	@Value("cat")
    String name;

	//..toString
}
```

##### 2.6.1.2 三种注入方式

1. 属性上直接标注

```java
@Data
public class Cat {
	@Autowired
    Person master;
```

2. 构造器注入

```java
@Autowired
public Cat(Person master) {
    this.master = master;
}
```

3. setter方法

```java
@Autowired
public void setMaster(Person master) {
    this.master = master;
}
```



启动类测试得,`Cat`里面已经依赖了`Person`

```sh
Cat{master=Person(name=思无邪), name='cat'}
```

##### 2.6.1.3 注入的bean不存在

如果查找的bean不存在，就会报错

```sh
Caused by: org.springframework.beans.factory.NoSuchBeanDefinitionException: No qualifying bean of type 'com.linkedbear.spring.basic_di.d_autowired.bean.Person' available: expected at least 1 bean which qualifies as autowire candidate. Dependency annotations: {}
```

如果不希望程序报错就注解上添加一个属性

```java
@Autowired(required = false) //默认为true
```

再次启动类测试得，所求依赖为null,不会报错。

```sh
Cat{master=null, name='cat'}
```

#### 2.6.2 @Autowired在配置类的使用

由于配置类上下文并没有`Person`的注册，也就没有`person()` 进行调用。就可以使用`@Autowired`进行自动注入**，但高版本可以不适用，也会自动注入**。

```java
@Configuration
@ComponentScan("org.example.AnnotationBean")
public class AnnotationConfiguration {

    @Bean
   //@Autowired
    public Cat cat(Person person) {
        Cat cat = new Cat();
        cat.setMaster(person);
        cat.setName("阿黄");
        return cat;
    }
}
```



#### 2.6.3 多个相同类型Bean的自动注入

创建多个bean对象

```java
@Data
@Component("admin")
@NoArgsConstructor
public class Person {

private String name = "思无邪";
}

//配置类中============

@Bean
public Person master() {
    Person person = new Person();
    person.setName("思无邪1");
    return person;
}
```

启动类测试得,发现两个`bean `,spring不知道该注入哪一个

```sh
 expected single matching bean but found 2: admin,master
```



##### 2.6.3.1 @Qualifier：指定注入Bean的名称

```java
public Cat cat(@Qualifier("master") Person person) {
    Cat cat = new Cat();
    cat.setMaster(person);
    cat.setName("阿黄");
    return cat;
}
//Cat{master=Person(name=思无邪1), name='cat'}
```

注意，`@Qualifier` 要在 `@Autowired` 下面使用

```java
@Autowired
@Qualifier("master")
private Person person;
```



##### 2.6.3.2 @Primary:默认Bean

`@Primary` 注解的使用目标是被注入的 Bean ，在一个应用中，一个类型的 Bean 注册只能有一个，它配合 `@Bean` 使用，可以指定默认注入的 Bean ：

```java
@Bean
@Primary
public Person master() {
    Person person = new Person();
    person.setName("思无邪1");
    return person;
}
```

`@Qualifier` 不受 `@Primary` 的干扰。

xml中可以指定`<bean>` 中的`primary`属性为true，与`@Primary` 注解一致



##### 2.6.3.3 其他方法

只要改变量名为对应的`bean`的id即可

```java
@Autowired
private Person master;
//Cat{master=Person(name=思无邪1), name='cat'}
```



##### 2.6.3.4 【面试题】 @Autowired注入的原理逻辑 

> **先拿属性对应的类型，去 IOC 容器中找 Bean ，如果找到了一个，直接返回；如果找到多个类型一样的 Bean ， 把属性名拿过去，跟这些 Bean 的 id 逐个对比，如果有一个相同的，直接返回；如果没有任何相同的 id 与要注入的属性名相同，则会抛出 `NoUniqueBeanDefinitionException` 异常。**



#### 2.6.4 多个相同类型的Bean的全部注入

**注入一个用单个对象接收，注入一组对象就用集合来接收**

~~~java
public class Cat {
@Autowired
private List<Person> persons;
    
//persons=[Person(name=思无邪), Person(name=思无邪1)]
~~~



#### 2.6.5 JSR250-@Resource

> JSR 全程 **Java Specification Requests** ，它定义了很多 Java 语言开发的规范，有专门的一个组织叫 JCP ( Java Community Process ) 来参与定制。
>
> 有关 JSR250 规范的说明文档可参考官方文档：[jcp.org/en/jsr/deta…](https://link.juejin.cn/?target=https%3A%2F%2Fjcp.org%2Fen%2Fjsr%2Fdetail%3Fid%3D250)

回到正题，`@Resource` 也是用来属性注入的注解，它与 `@Autowired` 的不同之处在于：**`@Autowired` 是按照类型注入，`@Resource` 是直接按照属性名 / Bean的名称注入**。

**`@Resource` 注解相当于标注 `@Autowired` 和 `@Qualifier`**

先导入依赖

```xml
<dependency>
    <groupId>javax.annotation</groupId>
    <artifactId>javax.annotation-api</artifactId>
    <version>1.3.2</version>
</dependency>
```

`@Resource` 到上下文查找bean。如果上下文没有就使用到对应的类(`@Component`注解的类)

```java
@Resource
private Person person;

@Bean
public Person master() {
    Person person = new Person();
    person.setName("思无邪123");
    return person;
}

//Cat{master=Person(name=思无邪123), name='cat'
```

#### 2.6.6 JSR330-@Inject

JSR330 也提出了跟 `@Autowired` 一样的策略，它也是**按照类型注入**。不过想要用 JSR330 的规范，需要额外导入一个依赖：

~~~xml
<!-- jsr330 -->
<dependency>
    <groupId>javax.inject</groupId>
    <artifactId>javax.inject</artifactId>
    <version>1</version>
</dependency>
~~~

使用与 SpringFramework 原生的 `@Autowired` + `@Qualifier` 一样了：

```java
@Component
public class Cat {
    
    @Inject // 等同于@Autowired
    @Named("admin") // 等同于@Qualifier
    private Person master;
```

##### 2.6.6.1 @Autowired与@Inject对比

查看包名

> **import** org.springframework.beans.factory.**annotation**.Autowired; 
>
> **import** javax.inject.Inject;

如果万一项目中没有 SpringFramework 了，那么 `@Autowired` 注解将失效，但 `@Inject` 属于 **JSR 规范，不会因为一个框架失效而失去它的意义**，只要导入其它支持 JSR330 的 IOC 框架，它依然能起作用。

#### 2.6.7 【面试题】依赖注入的注入方式



![image-20230413095848439](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/13/image-20230413095848439.png)



#### 2.6.8 【面试题】自动注入的注解对比

[面试题](#2.9.3 依赖注入具体是如何注入的？)

![image-20230413095926978](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/13/image-20230413095926978.png)



`@Qualifier` ：如果被标注的成员/方法在根据类型注入时发现有多个相同类型的 Bean ，则会根据该注解声明的 name 寻找特定的 bean

`@Primary` ：如果有多个相同类型的 Bean 同时注册到 IOC 容器中，使用 “根据类型注入” 的注解时会注入标注 `@Primary` 注解的 bean

---



### 2.7 复杂类型注入           

包括如下类型:

- 数组
- List / Set
- Map
- Properties



#### 2.7.1 构造复杂Person

```java
@Data
@Component
public class Person {

    private String name = "思无邪";
    private String[] names;
    private List<String> tels;
    private Set<Cat> cats;
    private Map<String, Object> events;
    private Properties props;
}
```



#### 2.7.2 xml复杂注入

```xml
<!-- 注意开启注解功能，否则cat的值不会导入 -->
<context:annotation-config></context:annotation-config>
<bean id="cat" class="org.example.AnnotationBean.Cat"/>

<bean class="org.example.AnnotationBean.Person">
<property name="names">
    <array>
        <value>张三</value>
        <value>李四</value>
    </array>
</property>

<property name="tels">
    <list>
        <value>1</value>
        <value>2</value>
    </list>
</property>
<property name="cats">
    <set>
        <ref bean="cat"/>
        <bean class="org.example.AnnotationBean.Cat"/>
    </set>
</property>

<property name="events">
    <map>
        <entry key="8:00" value="起床了"/>
        <entry key="9:00" value-ref="cat"/>
        <entry key="11:00">
            <bean class="org.example.AnnotationBean.Cat"/>
        </entry>
    </map>
</property>

<property name="props">
    <props>
        <prop key="sex">男</prop>
        <prop key="age">18</prop>
    </props>
</property>
</bean>
```

```sh
Person(name=思无邪, names=[张三, 李四], tels=[1, 2], cats=[Cat{name='cat'}], events={8:00=起床了, 9:00=Cat{name='cat'}, 11:00=Cat{name='cat'}}, props={sex=男, age=18})
```



#### 2.7.3 注解复杂注入

为了能演示 Bean 的引用，咱给 `Cat` 加上 `@Component` 注解，并带上名称：

```java
@Component("miaomiao")
public class Cat {
    private String name = "cat";
```

借助SpEL表达式进行注解注入



```java
package org.example.AnnotationBean;

import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.Set;

/**
 * @author wuxie
 * @date 2023/4/10 22:32
 * @description 该文件的描述 todo
 */
@Component
@Data
public class Person {

    private String name = "思无邪";

    @Value("#{new String[] {'张三', '张仨'}}")
    private String[] names;

    @Value("#{{'3333','333','33'}}")
    private List<String> tels;

    @Value("#{{@miaomiao,new org.example.AnnotationBean.Cat()}}")
    private Set<Cat> cats;

    @Value("#{{'喵喵':@miaomiao.name,'猫猫':new org.example.AnnotationBean.Cat()}}")
    private Map<String, Object> events;

    @Value("#{{'123':'你好','234':'我好'}}")
    private Properties props;
}

//Person(name=思无邪, names=[张三, 张仨], tels=[3333, 333, 33], cats=[Cat{name='cat'}], events={喵喵=cat, 猫猫=Cat{name='cat'}}, props={123=你好, 234=我好})
```



### 2.8 回调注入&延迟注入

#### 2.8.1 回调注入

##### 2.8.1.1 回调根源:Aware

回调注入的核心是一个叫 **`Aware`** 的接口，它来自 SpringFramework 3.1 ：

```java
public interface Aware {

}
```

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/985ddf54f86a42c8a9484310afdaa719~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

常用的几个回调接口

![image-20230414092443954](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230414092443954.png)

这里面大部分接口，其实在当下的 SpringFramework 5 版本中，借助 `@Autowired` 注解就可以实现注入了，根本不需要这些接口，只有最后面两个，是因 Bean 而异的，还是需要 **Aware** 接口来帮忙注入。

##### 2.8.1.2 ApplicationContextAware的使用

> 通过 ApplicationContextAware 接口可以让 Bean 类具有以下能力：
>
> 1. 获取 Spring 容器上下文（ApplicationContext）
> 2. 在 Bean 初始化的时候进行一些操作
> 3. 在 Bean 销毁前进行一些操作

###### 1.创建bean

~~~java
public class AwaredTestBean implements ApplicationContextAware {

    private ApplicationContext ctx;

    @Override
    public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
        this.ctx = applicationContext;
    }
}
~~~



在`AwaredTestBean` 注册初始化时，就调用`setApplicationContext` 将`ApplicationContext` 传给它，然后就可以使用`ApplicationContext`了

```java
public class AwaredTestBean implements ApplicationContextAware {

    private ApplicationContext ctx;

    public void printfBeanNames() {
        Stream.of(ctx.getBeanDefinitionNames()).forEach(System.out::println);
    }

    @Override
    public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
        this.ctx = applicationContext;
    }
}
```



###### 2.创建配置类

```java
@Configuration
public class AwareConfiguration {

    @Bean
    public AwaredTestBean awaredTestBean() {
        return new AwaredTestBean();
    }
}
```



###### 3.编写启动类

```java
public class AnnotationConfigApplication {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(AwareConfiguration.class);
        AwaredTestBean bean = ctx.getBean(AwaredTestBean.class);
        bean.printfBeanNames();
    }
}
/*
org.springframework.context.annotation.internalConfigurationAnnotationProcessor
org.springframework.context.annotation.internalAutowiredAnnotationProcessor
org.springframework.context.annotation.internalCommonAnnotationProcessor
org.springframework.context.event.internalEventListenerProcessor
org.springframework.context.event.internalEventListenerFactory
awareConfiguration
awaredTestBean
*/
```

由结果得，容器中的bean都被打印了，说明`ApplicationContext` 注入到`AwaredTestBean`成功。

##### 2.8.1.3 BeanNameAware的使用

> BeanNameAware 是 Spring 框架中的一个接口，它允许 Bean 对象在实例化和初始化时，获取它在 Spring 容器中注册的 Bean 的名称。

如果当前的 bean 需要依赖它本身的 name ，使用 `@Autowired` 就不好使了，这个时候就得使用 `BeanNameAware` 接口来辅助注入当前 bean 的 name 了。

###### 1.修改bean

```java
public class AwaredTestBean implements ApplicationContextAware, BeanNameAware {

    private ApplicationContext ctx;

    private String name;

    //......

    @Override
    public void setBeanName(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }
}
```



###### 2.修改启动类

```java
public class AnnotationConfigApplication {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(AwareConfiguration.class);
        AwaredTestBean bean = ctx.getBean(AwaredTestBean.class);
        bean.printfBeanNames();
        System.out.println("---------");
        System.out.println(bean.getName());
    }
}
/*
......
awareConfiguration
awaredTestBean
---------
awaredTestBean
*.
```



###### 3.NamedBean

其实，`BeanNameAware` 还有一个可选的搭配接口：**`NamedBean`** ，它专门提供了一个 `getBeanName` 方法，用于获取 bean 的 name 。

所以说，如果给上面的 `AwaredTestBean` 再实现 `NamedBean` 接口，那就不需要自己定义 `getName` 或者 `getBeanName` 方法，直接实现 `NamedBean` 定义好的 `getBeanName` 方法即可。



#### 2.8.2 延迟注入

[延迟查找](#7.3 优化-延迟查找 )

##### 2.8.2.1 setter的延迟注入

之前咱在写 setter 注入时，直接在 setter 中标注 `@Autowired` ，并注入对应的 bean 即可。如果使用延迟注入，则注入的就应该换成 `ObjectProvider` ：

```java
@Data
@Component
public class Cat {


    private Person person;

    @Autowired
    public void setPerson(ObjectProvider<Person> person) {
        // 有才注入
        this.person = person.getIfAvailable();
    }
```

##### 2.8.2.2 构造器的延迟注入

**与setter注入类似**

```java
@Data
@Component
public class Cat {


    private Person person;

    @Autowired
    public Cat(ObjectProvider<Person> person) {
        // 有才注入
        this.person = person.getIfAvailable();
    }
```



##### 2.8.2.2 构造器的延迟注入

属性直接注入是不能直接注入 Bean 的，只能注入 `ObjectProvider` ，通常也不会这么干，因为这样注入了之后，每次要用这个 Bean 的时候都得判断一次：

```java
@Autowired
private ObjectProvider<Person> people;

@Override
public String toString() {
    return "Cat{" +
        // 每次使用都要判断一次
            "name='" + people.getIfAvailable() + '\'' +
            '}';
}
```



#### 2.8.3 【面试题】依赖注入的注入方式

![image-20230414100119895](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230414100119895.png)



### 2.9 【面试题】 依赖注入

#### 2.9.1 依赖注入的目的和优点？

首先，依赖注入作为 IOC 的实现方式之一，目的就是**解耦**，我们不再需要直接去 new 那些依赖的类对象（直接依赖会导致对象的创建机制、初始化过程难以统一控制）；而且，如果组件存在多级依赖，依赖注入可以将这些依赖的关系简化，开发者只需要定义好谁依赖谁即可。

除此之外，依赖注入的另一个特点是依赖对象的**可配置**：通过 xml 或者注解声明，可以指定和调整组件注入的对象，借助 Java 的多态特性，可以不需要大批量的修改就完成依赖注入的对象替换（面向接口编程与依赖注入配合近乎完美）。

#### 2.9.2 谁把什么注入给谁了？

由于组件与组件之间的依赖只剩下成员属性 + 依赖注入的注解，而注入的注解又被 SpringFramework 支持，所以这个问题也好回答：**IOC 容器把需要依赖的对象注入给待注入的组件**。

#### 2.9.3 依赖注入具体是如何注入的？

[表格](#2.6.8 【面试题】自动注入的注解对比)

关于 `@Autowired` 注解的注入逻辑，在第 9 章 1.3.4 节有提过；`@Resource` 和 `@Inject` 的注入方式也都在第 9 章的 1.8 节罗列出来了，小伙伴们可以根据表格内容整理，灵活回答即可。

### 2.9.4 使用setter注入还是构造器注入？

这个问题，最好的保险回答是引用官方文档，而官方文档在不同的版本推荐的注入方式也不同，具体可参照如下回答：

- SpringFramework **4.0.2** 及之前是推荐 setter 注入，理由是**一个 Bean 有多个依赖时，构造器的参数列表会很长**；而且如果 **Bean 中依赖的属性不都是必需的话，注入会变得更麻烦**；
- **4.0.3** 及以后官方推荐构造器注入，理由是**构造器注入的依赖是不可变的、完全初始化好的，且可以保证不为 null** ；
- 当然 **4.0.3** 及以后的官方文档中也说了，如果**真的出现构造器参数列表过长的情况，可能是这个 Bean 承担的责任太多，应该考虑组件的责任拆解**。

![img](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/189e28f36ac94ad0b3948502df59e365~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

## 3.依赖查找与依赖注入的对比

- 作用目标不同

  - 依赖注入的作用目标通常是**类成员**

  - 依赖查找的作用目标可以是方法体内，也可以是方法体外

    - 在方法体内，依赖查找通常指的是组件主动从容器中获取其所需的依赖对象。在这种情况下，组件需要通过容器提供的API来获取依赖对象的引用，例如在Spring容器中可以使用ApplicationContext.getBean()方法进行依赖查找。

      而在方法体外，依赖查找通常指的是**容器在组装应用程序时自动完成对组件之间依赖关系的查找并注入依赖对象**。在这种情况下，组件只需要声明自己需要哪些依赖，容器就会负责将这些依赖注入到组件中，从而使得组件能够正常运行。

- 实现方式不同

  - 依赖注入通常借助一个**上下文被动的接收**
  - 依赖查找通常主动使用上下文搜索



## 4.IOC-Annotation

后面我们会经常注解进行声明式开发。前面使用 `ClassPathXmlApplicationContext` 对应类路径下的xml驱动。而我们开始使用注解配置驱动，就使用到了`AnnotationConfigApplicationContext`



### 4.1 注解驱动IOC的依赖查找

#### 4.1.1 配置类编写和Bean的注册

对于xml文件作为驱动，注解驱动需要的是**配置类**。只需要在类上面标注一个`@Configuration` 注解即可

> "Configurable" 和 "Configuration" 的区别在于前者是一个可配置的对象，而后者是一个用于定义 bean 的类或者方法。

```java
@Configuration
public class QuickStartConfiguration {
}
```

在xml里面，通过`<bean>` 标签进行声明Bean

配置类中，则使用`@Bean` 注解

```java
@Bean
public Person person() {
    return new Person();
}
```

解释：

1. 向IOC容器注册一个Person类型，id为person的Bean。返回值为注册类型，方法名代表Bean的id。也可以在注解上进行显式声明id,只不过属性为`name`

```java
@Bean(name = "myPerson")
public Person person() {
    return new Person();
}
```



#### 4.1.2 启动类初始化

 ```java
 public class AnnotationConfigApplication {
     public static void main(String[] args) {
         AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(QuickStartConfiguration.class);
         Person person = (Person) ctx.getBean("myPerson");
         Person bean = ctx.getBean(Person.class);
         System.out.println(person);
         System.out.println(bean);
     }
 }
 ```



### 4.2 注册驱动IOC的依赖注入

编码方式的依赖注入可以说是相当简单了，直接在创建对象后先别着急返回，把里面的值都 set 进去，再返回即可：

```java
@Bean(name = "myPerson")
public Person person() {
    Person person = new Person();
    person.setAge(20);
    person.setName("思无邪");
    return person;
}
```

```java
@Bean
public Cat cat() {
    Cat cat = new Cat();
    cat.setName("阿黄");
    // 直接使用上面的方法
    cat.setMaster(person());
    return cat;
}
```

相当于

```xml
<property name="name" value="test-cat"/>
<property name="master" ref="person"/>
```

启动测试

```java
Cat bean = ctx.getBean(Cat.class);
System.out.println(bean);
//Cat(master=Person(name=思无邪, age=20), name=阿黄)
```



### 4.3 组件注册和组件扫描

注册组件数量增多，那么`@Bean` 就很大工作量，所以spring框架开发出了用于快速注册组件的注解。**`模式注解 ( stereotype annotation )`**



#### 4.3.1 一切组件注册的根源：@Component

在类上标注 `@Component` 注解，即代表**该类**会被注册到 IOC 容器中作为一个 Bean 。

```java
@Data
@Component
public class Person {
	String name;
    int age;
}
```



如果想要指定Bean的名称，可以直接`@Component` 声明**value** 

```java
@Component(value = "myPerson")
@Component("myPerson")
```

如果不指定Bean的名称，默认规则为**“类名的首字母小写”**（Person类默认名称为`person`)



#### 4.3.2 组件扫描

配置类中`@Configuration`注解，感知不到`@Component`存在，会报错`NoSuchBeanDefinitionException`

所以使用`@ComponentScan`.

##### 4.3.2.1 @ComponentScan

在配置类上额外标注一个`@ComponentScan`,并指定要扫描的路径，会**扫描指定路径包及子包下的所有`@Component`组件**

```java
@Configuration
@ComponentScan("org.example.bean")
//@ComponentScan(basePackages = "org.example.bean")
```

如果不指定扫描路径，则**默认扫描本类所在包及子包下的所有 `@Component` 组件**。



##### 4.3.2.2 不适用@ComponentScan

```java
AnnotationConfigApplicationContext ctx = new 
// 通过路径进行
AnnotationConfigApplicationContext("org.example.bean");
Person bean = ctx.getBean(Person.class);
System.out.println(bean);
```



##### 4.3.2.3 xml启动组件扫描

```xml
<context:component-scan base-package="org.example.bean"/>
```

之后使用 `ClassPathXmlApplicationContext` 驱动，也是可以获取到 `Person` 的。



#### 4.3.3 组件注册的其他注解

> SpringFramework 为了迎合咱在进行 Web 开发时的三层架构，它额外提供了三个注解：`@Controller` 、`@Service` 、`@Repository` ，分别代表表现层、业务层、持久层。这三个注解的作用与 `@Component` 完全一致，其实它们的底层也就是 `@Component` ：

```java
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Component
public @interface Controller { ... }
```



4.3.4 `@Configuration也是@Component`

将包扫描范围扩大

```java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext("org.example");
String[] beanDefinitionNames = ctx.getBeanDefinitionNames();
Stream.of(beanDefinitionNames).forEach(System.out::println);

/*
componentScanConfiguration
quickStartConfiguration
...
*/
```

可见配置类也被注册到了IOC容器上

查看其源码得

```java
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Component
public @interface Configuration {
```

### 4.4 注解驱动和xml驱动互通

#### 4.4.1 xml 引入注解

需要开启注解配置，注册对应的配置类

```xml
<context:annotation-config/>
<bean class="org.example.config.ComponentScanConfiguration"/>
```

#### 4.4.2 注解引入xml

在注解配置中引入 xml ，需要在配置类上标注 `@ImportResource` 注解，并声明配置文件的路径：

```java
@Configuration
@ImportResource(value = "classpath:application.xml")
public class AnnotationConfiguration {
}

```



## 5.Bean常见的几种类型和作用域

### 5.1 Bean的类型

在 SpringFramework 中，对于 Bean 的类型，一般有两种设计：**普通 Bean 、工厂 Bean** 。以下分述这两种类型。



#### 5.1.1 普通Bean

之前创建的Bean都是普通的Bean

```java
@Component
public class Child {

}
```

#### 5.1.2 FactoryBean

SpringFramework 考虑到一些特殊的设计：Bean 的创建需要指定一些策略，或者依赖特殊的场景来分别创建，也或者一个对象的创建过程太复杂，使用 xml 或者注解声明也比较复杂。这种情况下，如果还是使用普通的创建 Bean 方式，以咱现有的认知就搞不定了。于是，SpringFramework 在一开始就帮我们想了办法，可以借助 **`FactoryBean`** 来使用工厂方法创建对象。



##### 5.1.2.1 介绍

`FactoryBean` 本身是一个接口，它本身就是一个创建对象的工厂。如果 Bean 实现了 `FactoryBean` 接口，则它本身将不再是一个普通的 Bean ，**不会在实际的业务逻辑中起作用，而是由创建的对象来起作用**。

```java
public interface FactoryBean<T> {
    // 返回创建的对象
    //指示某个参数、字段或方法的返回值可以为 null
    @Nullable
    T getObject() throws Exception;

    // 返回创建的对象的类型（即泛型类型）
    @Nullable
    Class<?> getObjectType();

    // 创建的对象是单实例Bean还是原型Bean，默认单实例
    default boolean isSingleton() {
        return true;
    }
}
```



#### 5.1.3 Factory的使用

##### 5.1.3.1 创建消费者

```java
public class Child {
    // 当前的小孩子想玩球
    private String wantToy = "ball";

    public String getWantToy() {
        return wantToy;
    }
}
```



创建几个使用对象

```java
public abstract class Toy {
    
    private String name;
    
    public Toy(String name) {
        this.name = name;
    }
    @Override
    public String toString() {
        return "Toy{" +
                "name='" + name + '\'' +
                '}';
    }
}
```



```java
public class Ball extends Toy { // 球
    
    public Ball(String name) {
        super(name);
    }
}
```



```java
package org.example.Bean_Type;

public class Car extends Toy { // 玩具汽车
    
    public Car(String name) {
        super(name);
    }
}
```



##### 5.1.3.2 创建玩具工厂

创建一个 `ToyFactoryBean` ，让它实现 `FactoryBean` 接口：

```java
public class ToyFactoryBean implements FactoryBean<Toy> {
    @Override
    public Toy getObject() throws Exception {
        return null;
    }

    @Override
    public Class<?> getObjectType() {
        return Toy.class;
    }
}
```

让它根据小孩子想要玩的玩具来决定生产哪种玩具，那咱就得在这里面注入 `Child` 。由于咱这里面使用的不是注解式自动注入，那咱就用 setter 注入吧：

```java
public class ToyFactoryBean implements FactoryBean<Toy> {

    private Child child;

    public void setChild(Child child) {
        this.child = child;
    }

    @Override
    public Toy getObject() throws Exception {
        switch (child.getWantToy()) {
            case "ball":
                return new Ball("ball");
            case "car":
                return new Car("car");
            default:
                return null;
        }
    }

    @Override
    public Class<?> getObjectType() {
        return Toy.class;
    }
}
```



##### 5.1.3.3 注册工厂类

**xml方式**

```xml
    <bean class="org.example.Bean_Type.Child" id="child"/>
    <bean class="org.example.Bean_Type.ToyFactoryBean" id="toyFactoryBean">
        <property name="child" ref="child"/>
    </bean>
```

**注解方式**

```java
@Configuration
public class BeanTypeConfiguration {
    @Bean
    public Child child() {
        return new Child();
    }

    @Bean
    public ToyFactoryBean toyFactoryBean() {
        ToyFactoryBean toyFactoryBean = new ToyFactoryBean();
        toyFactoryBean.setChild(child());
        return toyFactoryBean;
    }
}
```



##### 5.1.3.4 测试类

```java
public class AnnotationConfigApplication {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(BeanTypeConfiguration.class);
        Toy bean = ctx.getBean(Toy.class);
        //Toy{name='ball'}
        System.out.println(bean);
    }
}
```



#### 5.1.4 FactoryBean与Bean同时存在

修改配置文件 / 配置类，向 IOC 容器预先的创建一个 `Ball` ，这样 `FactoryBean` 再创建一个，IOC 容器里就会同时存在两个 `Toy` 了：

```java
@Bean
public Toy ball() {
    return new Ball("ball");
}

@Bean
public ToyFactoryBean toyFactoryBean() {
    ToyFactoryBean toyFactoryBean = new ToyFactoryBean();
    toyFactoryBean.setChild(child());
    return toyFactoryBean;
}
```

再次运行主类可以发现抛出了`NoUniqueBeanDefinitionException` 异常，`expected single matching bean but found 2: ball,toyFactoryBean`

提示有两个`Toy`,说明`FactoryBean` 注入了IOC容器。打印一下

```java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(BeanTypeConfiguration.class);
Map<String, Toy> beansOfType = ctx.getBeansOfType(Toy.class);
beansOfType.forEach((name, toy) -> {
    System.out.println(name + ":" + toy);
});
/*
ball:Toy{name='ball'}
toyFactoryBean:Toy{name='ball'}
*/
```

#### 5.1.5 FactoryBean创建Bean的时机

`ApplicationContext` 初始化 Bean 的时机默认是容器加载时就已经创建，那 `FactoryBean` 创建 Bean 的时机又是什么呢？

##### 5.1.5.1 FactoryBean的加载时机

给 `Toy` 的构造方法中添加一个控制台打印：

```java
public Toy(String name) {
    System.out.println("生成了一个" + name);
    this.name = name;
}
```

`ToyFactoryBean` 构造器添加打印信息

```java
public ToyFactoryBean() {
    System.out.println("ToyFactoryBean 初始化了...");
}
```

配置类

```java
@Configuration
public class BeanTypeConfiguration {
    @Bean
    public Child child() {
        return new Child();
    }

//    @Bean
//    public Toy ball() {
//        return new Ball("ball");
//    }

    @Bean
    public ToyFactoryBean toyFactoryBean() {
        ToyFactoryBean toyFactoryBean = new ToyFactoryBean();
        toyFactoryBean.setChild(child());
        return toyFactoryBean;
    }
}
```



测试启动类

```java
ApplicationContext ctx = new AnnotationConfigApplicationContext(BeanTypeConfiguration.class);
//ToyFactoryBean 初始化了...
```

只有 `ToyFactoryBean` 被初始化，说明 **`FactoryBean` 本身的加载是伴随 IOC 容器的初始化时机一起的**。

##### 5.1.5.2 FactoryBean的加载时机

```java
ApplicationContext ctx = new AnnotationConfigApplicationContext(BeanTypeConfiguration.class);
Toy bean = ctx.getBean(Toy.class);
/*
ToyFactoryBean 初始化了...
生成了一个ball
*/
```

也就得出：**`FactoryBean` 生产 Bean 的机制是延迟生产**。



#### 5.1.6 FactoryBean创建Bean的实例数

查看`FactoryBean`的源码可见一个`isSingleton` 方法,默认为true，代表默认是单实例的

```java
default boolean isSingleton() {
    return true;
}
```

修改启动类，连续取出两次`Toy`

```java
public class AnnotationConfigApplication {
    public static void main(String[] args) {
        ApplicationContext ctx = new AnnotationConfigApplicationContext(BeanTypeConfiguration.class);
        Toy bean = ctx.getBean(Toy.class);
        Toy bean1 = ctx.getBean(Toy.class);
        System.out.println(bean1 == bean);
    }
}
/*
ToyFactoryBean 初始化了...
生成了一个ball
true
*/
```



#### 5.1.7 取出FactoryBean本体

咱刚才一直都是拿 `Toy` 本体去取，取到的都是 `FactoryBean` 生产的 Bean 。一般情况下咱也用不到 `FactoryBean` 本体，但如果真的需要取，使用的方法也很简单：要么直接传 `FactoryBean` 的 class （很容易理解），也可以传 ID 。不过，**如果真的靠传 ID 的话，传配置文件 / 配置类声明的 ID 就不好使了，因为那样只会取出生产出来的 Bean** ：

```java
ApplicationContext ctx = new AnnotationConfigApplicationContext(BeanTypeConfiguration.class);
System.out.println(ctx.getBean("toyFactoryBean"));

//Toy{name='ball'}
```

取本体需要在Bean的id前面加"&"符号

```java
System.out.println(ctx.getBean("&toyFactoryBean"));
//org.example.Bean_Type.ToyFactoryBean@4f2b503c
```

#### 5.1.8 【面试题】BeanFactory与FactoryBean的区别

`BeanFactory` ：SpringFramework 中实现 IOC 的最底层容器（此处的回答可以从两种角度出发：从类的继承结构上看，它是最顶级的接口，也就是最顶层的容器实现；从类的组合结构上看，它则是最深层次的容器，`ApplicationContext` 在最底层组合了 `BeanFactory` ）

![image-20230416104607675](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230416104607675.png)

`FactoryBean` ：创建对象的工厂 Bean ，可以使用它来直接创建一些初始化流程比较复杂的对象



### 5.2 Bean的作用域

```java
package org.example.Bean_Type;

public class ScopeReviewDemo {
    // 类级别成员
    private static String classVariable = "";
    
    // 对象级别成员
    private String objectVariable = "";
    
    public static void main(String[] args) throws Exception {
        // 方法级别成员
        String methodVariable = "";
        for (int i = 0; i < args.length; i++) {
            // 循环体局部成员
            String partVariable = args[i];
            
            // 此处能访问哪些变量？
            /*
            classVariable、objectVariable 和 methodVariable,partVariable
             */
        }
        
        // 此处能访问哪些变量？
        //classVariable
    }
    
    public void test() {
        // 此处能访问哪些变量？
        //classVariable,objectVariable
    }
    
    public static void staticTest() {
        // 此处能访问哪些变量？
        // classVariable
    }
}
```



### 5.3 Spring内置的作用域

SpringFramework 中内置了 6 种作用域（5.x 版本）：

| 作用域类型  | 概述                                         |
| ----------- | -------------------------------------------- |
| singleton   | 一个 IOC 容器中只有一个【默认值】            |
| prototype   | 每次获取创建一个                             |
| request     | 一次请求创建一个（仅Web应用可用）            |
| session     | 一个会话创建一个（仅Web应用可用）            |
| application | 一个 Web 应用创建一个（仅Web应用可用）       |
| websocket   | 一个 WebSocket 会话创建一个（仅Web应用可用） |



#### 5.3.1 singleton 单实例Bean

![img](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/f1a1548eb64b49c797bc0155c43512bc~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

左边的几个定义的 Bean 同时引用了右边的同一个 `accountDao` ，对于这个 `accountDao` 就是单实例 Bean 。

SpringFramework 中默认所有的 Bean 都是单实例的，即：**一个 IOC 容器中只有一个**。下面咱演示一下单实例 Bean 的效果：

##### 5.3.1.1 创建Bean + 配置类

```java
public class Child {
    
    private Toy toy;
    
    public void setToy(Toy toy) {
        this.toy = toy;
    
     @Override
    public String toString() {
        return "Child{" +
                "toy=" + toy +
                '}';
    }
}
```

```java
// Toy 中标注@Component注解
@Component
public class Toy {
    
}
```

创建配置类，同时注册两个`Child`

 ```java
 @Configuration
 @ComponentScan("com.linkedbear.spring.bean.b_scope.bean")
 public class BeanScopeConfiguration {
     
     @Bean
     public Child child1(Toy toy) {
         Child child = new Child();
         child.setToy(toy);
         return child;
     }
     
     @Bean
     public Child child2(Toy toy) {
         Child child = new Child();
         child.setToy(toy);
         return child;
     }
     
 }
 ```

启动主类，打印child里面的toy 

~~~java
public class BeanScopeAnnoApplication {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(BeanScopeConfiguration.class);
        ctx.getBeansOfType(Child.class).forEach((name, child) -> {
            System.out.println(name + ":" + child);
        });
        
/*
child1:Child{toy=org.example.b_scope.Toy@11392934}
child2:Child{toy=org.example.b_scope.Toy@11392934}
*
/
    }
}
~~~

由此可见，child里面注入的toy是同一个



#### 5.3.2 prototype: 原型Bean

Spring 官方的定义是：**每次对原型 Bean 提出请求时，都会创建一个新的 Bean 实例。** 这里面提到的 ”提出请求“ ，包括任何依赖查找、依赖注入的动作，都算做一次 ”提出请求“ 。由此咱也可以总结一点：如果连续 `getBean()` 两次，那就应该创建两个不同的 Bean 实例；向两个不同的 Bean 中注入两次，也应该注入两个不同的 Bean 实例。

![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/c497ccbeac5e49eba21839c8c2f08a1a~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

> 其实对于**原型**这个概念，在设计模式中也是有对应的：**原型模式**。原型模式实质上是使用对象深克隆，乍看上去跟 SpringFramework 的原型 Bean 没什么区别，但咱仔细想，每一次生成的原型 Bean 本质上都还是一样的，只是可能带一些特殊的状态等等，



##### 5.3.2.1 修改Bean

```java
@Component
@Scope("prototype")
public class Toy {
    
}
```

作用域在`ConfigurableBeanFactory` 里面设置了

```java
public interface ConfigurableBeanFactory extends HierarchicalBeanFactory, SingletonBeanRegistry {

/**
 * Scope identifier for the standard singleton scope: {@value}.
 * <p>Custom scopes can be added via {@code registerScope}.
 * @see #registerScope
 */
String SCOPE_SINGLETON = "singleton";

/**
 * Scope identifier for the standard prototype scope: {@value}.
 * <p>Custom scopes can be added via {@code registerScope}.
 * @see #registerScope
 */
String SCOPE_PROTOTYPE = "prototype";
```



##### 5.3.2.2 测试运行

发现注入的两个toy都不是同一个

```java
child1:Child{toy=org.example.b_scope.Toy@6e9175d8}
child2:Child{toy=org.example.b_scope.Toy@7d0b7e3c}
```



##### 5.3.2.3 原型Bean的创建时机

```java
@Component
@Scope("prototype")
public class Toy {
    public Toy() {
        System.out.println("Toy create ...");
    }
}
```

测试运行

~~~java
public class BeanScopeAnnoApplication {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(BeanScopeConfiguration.class);
    }
}

/*
Toy create ...
Toy create ...
*/
~~~



#### 5.3.3 Web应用和作用域们

- request ：请求Bean，每次**客户端向 Web 应用服务器发起一次请求**，Web 服务器接收到请求后，由 SpringFramework 生成一个 Bean ，直到请求结束
- session ：会话Bean，每个客户端在与 Web 应用服务器**发起会话后**，SpringFramework 会为之生成一个 Bean ，直到会话过期
- application ：应用Bean，每个 **Web 应用在启动时**，SpringFramework 会生成一个 Bean ，直到应用停止（有的也叫 global-session ）
- websocket ：WebSocket Bean ，每个客户端在与 **Web 应用服务器建立 WebSocket 长连接时**，SpringFramework 会为之生成一个 Bean ，直到断开连接

上面 3 种可能小伙伴们还熟悉，最后一种 WebSocket 可能有些小伙伴还不了解，不了解没关系，这玩意用的也少，等回头小伙伴学习了 WebSocket 之后自己试一下就可以了，小册这里也不多展开了。





## 6.Bean的实例方式

**所有实例化指调用构造方法，创建新的对象**，**初始化指创建好新的对象后的属性赋值、组件注入等后续动作**



### 6.1 普通Bean的实例化

前面使用的`@Bean` 和`<bean>` 方式都是普通Bean的实例化，默认是单实例。IOC容器初始化的时候就已经被初始化

### 6.2 借助FactoryBean创建Bean

前面也写过

~~~java
public abstract class Toy {

    private String name;

    public Toy(String name) {
        System.out.println("生成了一个" + name);
        this.name = name;
    }

    @Override
    public String toString() {
        return "Toy{" +
                "name='" + name + '\'' +
                '}';
    }
}
~~~

~~~java
blic class ToyFactoryBean implements FactoryBean<Toy> {

    private Child child;

    public ToyFactoryBean() {
        System.out.println("ToyFactoryBean 初始化了...");
    }

    public void setChild(Child child) {
        this.child = child;
    }

    @Override
    public Toy getObject() throws Exception {
        switch (child.getWantToy()) {
            case "ball":
                return new Ball("ball");
            case "car":
                return new Car("car");
            default:
                return null;
        }
    }

    @Override
    public Class<?> getObjectType() {
        return Toy.class;
    }

~~~

配置类进行注册

~~~java
@Bean
public ToyFactoryBean toyFactoryBean() {
    ToyFactoryBean toyFactoryBean = new ToyFactoryBean();
    toyFactoryBean.setChild(child());
    return toyFactoryBean;
}
~~~

只要注册工厂，IOC会自动识别，并且默认在第一次获取时创建对应的Bean并缓存（针对默认的单实例）

### 6.3 借助静态工厂创建Bean

#### 6.3.1 创建Bean和工厂

~~~java
public class Car {
    
    public Car() {
        System.out.println("Car constructor run ...");
    }
}
~~~

创建一个静态工厂

~~~java
public class CarStaticFactory {
    
    public static Car getCar() {
        return new Car();
    }
}
~~~

#### 6.3.2 xml注册

静态工厂的使用通常运用于 xml 方式比较多（主要是**注解驱动没有直接能让它起作用的注解**，编程式配置又可以直接调用，显得没那么大必要，下面会演示），咱下面创建一个 `bean-instantiate.xml` 文件，在这里面编写关于静态工厂的使用方法：

~~~xml
<bean id="car1" class="com.linkedbear.spring.bean.c_instantiate.bean.Car"/>

<bean id="car2" class="com.linkedbear.spring.bean.c_instantiate.bean.CarStaticFactory" factory-method="getCar"/>
~~~

上面的普通的Bean注册，下面是直接引用静态工厂，声明工厂方法`factory-method`

```java
ClassPathXmlApplicationContext ctx = new ClassPathXmlApplicationContext("testBean.xml");
ctx.getBeansOfType(Car.class).forEach((name, type) -> {
    System.out.println(name + ":" + type);
});
/*
car1:org.example.Bean_Type.Car@7a8c8dcf
car2:org.example.Bean_Type.Car@24269709

*/
```



#### 6.3.3 测试静态工厂会在IOC容器吗

```java
System.out.println(ctx.getBean(BeanFactory.class));
//factory.NoSuchBeanDefinitionException
```

可见抛出了异常，得出结论:**静态工厂本身不会被注册IOC容器中**



#### 6.3.4 编程式使用静态工厂

并没有提供关于静态工厂相关的注解，使用注解配置类+编程式使用静态工厂。

在配置类里面进行注册

```java
@Bean
public Car car2() {
    return CarStaticFactory.getCar();
}
```



### 6.4 借助实例工厂创建Bean

类似于静态工厂

#### 6.4.1 创建实例工厂

与静态工厂的区别在于，不是**static** 方法

~~~java
public class CarInstanceFactory {
    
    public Car getCar() {
        return new Car();
    }
}
~~~

#### 6.4.2 配置xml

~~~xml
<bean id="car3" class="org.example.Bean_Type.CarInstanceFactory"/>
<bean id="car4" factory-bean="car3" factory-method="getCar"/>
~~~

先注册实例工厂，然后使用`factory-bean` 和 `factory-method` 属性也可以完成 Bean 的创建

#### 6.4.3 测试运行

~~~java
car1:org.example.Bean_Type.Car@24269709
car2:org.example.Bean_Type.Car@2aceadd4
car4:org.example.Bean_Type.Car@24aed80c
org.example.Bean_Type.CarInstanceFactory@2925bf5b
~~~

并且实例工厂也是被注册到IOC容器

#### 6.4.4 编程式使用实例工厂

类似静态工厂只不过要传入工厂对象

```java
@Bean
public Car car3(CarInstanceFactory carInstanceFactory) {
    return carInstanceFactory.getCar();
}
```



# 4.spring相关的信息

## 1.生命周期

### 1.1生命周期阶段

![img](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/4376c21f525e4c11b4fc07148e65e8e0~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

一个对象从被创建，到被垃圾回收，可以宏观的划分为 5 个阶段：

- **创建 / 实例化阶段**：此时会调用类的构造方法，产生一个新的对象
- **初始化阶段**：此时对象已经创建好，但还没有被正式使用，可能这里面需要做一些额外的操作（如预初始化数据库的连接池）
- **运行使用期**：此时对象已经完全初始化好，程序正常运行，对象被使用
- **销毁阶段**：此时对象准备被销毁，已不再使用，需要预先的把自身占用的资源等处理好（如关闭、释放数据库连接）
- **回收阶段**：此时对象已经完全没有被引用了，被垃圾回收器回收

### 1.2 spring能干预的生命周期阶段

因为spring的控制反转功能，Bean的创建和回收都是交给IOC容器，我们无法进行操作。而运行使用期Bean已被初始化完我们只能使用，无法进行额外设计了。所以只剩下**初始化和销毁的两个阶段** 可以操作。



如果进行干预Bean的初始化和销毁呢。

由**Servlet**得，有两个方法`init` 和 `destory` ,都不是自己调用，而是由**Web容器（Tomcat等）进行调用**。设计思想是**回调机制**，并非自己设计，而是父类，第三方接口设计，然后用**容器，框架等**来调用。与前面spring的`Aware` 接口调用思想一样。

> 生命周期的触发，更适合叫回调，因为生命周期方法是咱定义的，但方法被调用，是框架内部帮我们调的，那也就可以称之为 “回调” 了。

## 2. init-method&destroy-method

### 2.1 创建Bean

```java
public class Cat {
    
    private String name;
    
    public void setName(String name) {
        this.name = name;
    }
    
    public void init() {
        System.out.println(name + "被初始化了。。。");
    }
    public void destroy() {
        System.out.println(name + "被销毁了。。。");
    }
}
```

### 2.2 创建配置类/配置文件

`init-method` 设置初始化方法，`destroy-method` 销毁方法

~~~xml
<bean id="cat" class="org.example.BeanLife.Cat" init-method="init" destroy-method="destroy"/>
~~~

创建配置类

~~~java
@Configuration
public class BeanConfiguration {

    @Bean(initMethod = "init",destroyMethod = "destroy")
    public Cat cat() {
        Cat cat = new Cat();
        cat.setName("ahuang");
        return cat;
    }
}
~~~

初始化和销毁方法的要求特征

注意一点，这些配置的初始化和销毁方法必须具有以下特征：（原因一并解释）

- 方法访问**权限无限制要求**（ SpringFramework 底层会反射调用的）
- 方法**无参数**（如果真的设置了参数，SpringFramework 也不知道传什么进去）
- 方法**无返回值**（返回给 SpringFramework 也没有意义）
- 可以抛出异常（异常不由自己处理，交予 SpringFramework 可以打断 Bean 的初始化 / 销毁步骤）

### 2.3 测试效果

```java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(BeanConfiguration.class);
System.out.println(ctx.getBean(Cat.class));
ctx.close();
/*
ahuang被初始化了。。。
org.example.BeanLife.Cat@5562c41e
ahuang被销毁了。。。
*/
```

由此可以得出结论：**在 IOC 容器初始化之前，默认情况下 Bean 已经创建好了，而且完成了初始化动作；容器调用销毁动作时，先销毁所有 Bean ，最后 IOC 容器全部销毁完成。**

### 2.4 Bean的初始化流程顺序

```java
public class Cat {
    
    private String name;

    public Cat() {
        System.out.println("构造方法");
    }

    public void setName(String name) {
        System.out.println("setName方法执行了。。。");
        this.name = name;
    }
    
    public void init() {
        System.out.println(name + "被初始化了。。。");
    }
    public void destroy() {
        System.out.println(name + "被销毁了。。。");
    }
}
```

重新启动类进行测试:

```sh
构造方法
setName方法执行了。。。
ahuang被初始化了。。。
org.example.BeanLife.Cat@53fdffa1
ahuang被销毁了。。。
```

**Bean 的生命周期中，是先对属性赋值，后执行 `init-method` 标记的方法**。



## 3.JSR250规范

上面的方法，都是咱手动声明注册的 Bean ，对于那些使用模式注解的 Bean ，这种方式就不好使了，因为**没有可以让你声明 `init-method` 和 `destroy-method` 的地方了，`@Component` 注解上也只有一个 `value` 属性而已**。这个时候咱就需要学习一种新的方式，这种方式专门配合注解式注册 Bean 以完成全注解驱动开发，那就是如标题所说的 **JSR250 规范**。

除了像`@Resource` 这样的自动注入注解，还有负责生命周期的注解。

**`@PostConstruct`** 、**`@PreDestroy`** 两个注解，分别对应 `init-method` 和 `destroy-method` 。

### 3.1 创建Bean

~~~java
@Component
public class Pen {

    private Integer ink;


    @PostConstruct
    // 可以是private
    private void addInk() {
        System.out.println("添加墨水");
        this.ink = 100;
    }

    @PreDestroy
    public void outwellInk() {
        System.out.println("清空墨水");
        this.ink = 0;
    }

    @Override
    public String toString() {
        return "Pen{" + "ink=" + ink + '}';
    }
}
~~~

被 `@PostConstruct` 和 `@PreDestroy` 注解标注的方法，与 `init-method` / `destroy-method` 方法的声明要求是一样的，访问修饰符也可以是 private 。

~~~java
import javax.annotation.PostConstruct;
import javax.annotation.PreDestroy;
import javax.annotation.Resource;
import org.springframework.stereotype.Component;

AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext("org.example.BeanLife");
ctx.close();
/*
添加墨水
清空墨水
*/
~~~



### 3.2 JSR250 规范与init-method共存

如果不使用 `@Component` 注解来注册 Bean 而转用 `<bean>` / `@Bean` 的方式，那 `@PostConstruct` 与 `@PreDestroy` 注解是可以与 `init-method` / `destroy-method` 共存的

~~~java
//@Component
public class Pen {

    private Integer ink;

    public void open() {
        System.out.println("打开钢笔");
    }

    public void close() {
        System.out.println("关闭钢笔");
    }

    @PostConstruct
    private void addInk() {
        System.out.println("添加墨水");
        this.ink = 100;
    }

    @PreDestroy
    public void outwellInk() {
        System.out.println("清空墨水");
        this.ink = 0;
    }
}
~~~

修改配置类

~~~java
@Bean(initMethod = "open", destroyMethod = "close")
public Pen pen() {
    return new Pen();
}
~~~

启动主类测试

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(BeanConfiguration.class);
ctx.close();
/*
添加墨水
打开钢笔
清空墨水
关闭钢笔
*/	
~~~

并且可见**JSR250规范优先级大于init/destory**

## 4. InitializingBean&DisposableBean

这两个家伙实际上是两个接口，而且是 SpringFramework 内部预先定义好的两个关于生命周期的接口。他们的触发时机与上面的 `init-method` / `destroy-method` 以及 JSR250 规范的两个注解一样，都是在 Bean 的初始化和销毁阶段要回调的



### 4.1 创建Bean

实现接口后要实现对应的方法

~~~java
@Component
public class Person implements InitializingBean, DisposableBean {
    @Override
    public void destroy() throws Exception {
		
    }

    @Override
    public void afterPropertiesSet() throws Exception{

    }
}

~~~



### 4.2 三种生命周期定义方式并存

~~~java
public void open() {
        System.out.println("init-method - 打开钢笔。。。");
    }
    
    public void close() {
        System.out.println("destroy-method - 合上钢笔。。。");
    }
    
    @PostConstruct
    public void addInk() {
        System.out.println("@PostConstruct - 钢笔中已加满墨水。。。");
        this.ink = 100;
    }
    
    @PreDestroy
    public void outwellInk() {
        System.out.println("@PreDestroy - 钢笔中的墨水都放干净了。。。");
        this.ink = 0;
    }
    
    @Override
    public void afterPropertiesSet() throws Exception {
        System.out.println("InitializingBean - 准备写字。。。");
    }
    
    @Override
    public void destroy() throws Exception {
        System.out.println("DisposableBean - 写完字了。。。");
    }
~~~

注解驱动方式注入

~~~java
@Bean(initMethod = "open", destroyMethod = "close")
public Pen3 pen() {
    return new Pen3();
}
~~~

测试得三种方式的优先级

**`@PostConstruct` → `InitializingBean` → `init-method`** 。

## 5.原型Bean的生命周期

对于原型 Bean 的生命周期，使用的方式跟上面是完全一致的，只是它的触发时机就不像单实例 Bean 那样了。

单实例 Bean 的生命周期是陪着 IOC 容器一起的，容器初始化，单实例 Bean 也跟着初始化（当然不绝对，后面会介绍延迟 Bean ）；容器销毁，单实例 Bean 也跟着销毁。**原型 Bean 由于每次都是取的时候才产生一个，所以它的生命周期与 IOC 容器无关**。



### 5.1 创建Bean和配置类

~~~java
@Configuration
public class PrototypeLifecycleConfiguration {

@Bean(initMethod = "open", destroyMethod = "close")
@Scope(ConfigurableBeanFactory.SCOPE_PROTOTYPE)
public Pen pen() {
    return new Pen();
}
}
~~~

测试得并不会执行open方法,所以**原型Bean的创建不随IOC的初始化而创建**



### 5.2 获得Bean时初始化

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(PrototypeLifecycleConfiguration.class);
System.out.println(ctx.getBean(Pen.class));
ctx.close();

/*
添加墨水
打开钢笔
Pen{ink=100}
*/
~~~

所以初始化动作与单实例Bean完全一致,就是时间不一样;

这里有个问题我`ctx.close`,但是闭关没有输出`Destory` 相关的信息,因为前面说了原型Bean的生命周期与IOC容器无关.



### 5.3 原型Bean的销毁

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(PrototypeLifecycleConfiguration.class);
Pen bean = ctx.getBean(Pen.class);
System.out.println(bean);
//        ctx.close();
ctx.getBeanFactory().destroyBean(bean);

/*
添加墨水
hello open pen
打开钢笔
Pen{ink=100}
清空墨水
bye close pen
*/
~~~

可以得知执行了 `@PreDestroy` 注解的方法以及`InitializingBean, DisposableBean的方法`,但没有执行`destroyMethod` 的执行

结论:原型Bean的销毁,不处理`destroyMethod` 标注的方法



## 6.【面试题】 控制Bean生命周期的三种方法

|            | init-method & destroy-method              | @PostConstruct & @PreDestroy    | InitializingBean & DisposableBean |
| ---------- | ----------------------------------------- | ------------------------------- | --------------------------------- |
| 执行顺序   | 最后                                      | 最先                            | 中间                              |
| 组件耦合度 | 无侵入（只在 `<bean>` 和 `@Bean` 中使用） | 与 JSR 规范耦合                 | 与 SpringFramework 耦合           |
| 容器支持   | xml 、注解原生支持                        | 注解原生支持，xml需开启注解驱动 | xml 、注解原生支持                |
| 单实例Bean | √                                         | √                               | √                                 |
| 原型Bean   | 只支持 init-method                        | √                               | √                                 |



## 7.IOC容器的详细对比-BeanFactory

### 1.BeanFactory和它的接口

![img](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/c3fb460e3a9343d5b47a101d9bc5a9df~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

#### 1.1 BeanFactory是根容器【掌握】

`BeanFactory` 作为 SpringFramework 中最顶级的容器接口，它的作用一定是最简单、最核心的。下面咱先来看一看文档注释 ( javadoc ) 中的描述。

> 用于访问Spring bean容器的根接口。
> 这是bean容器的基本客户端视图；ListableBeanFactory和org.springframework.beans.factory.config.ConfigurableBeanFactory等其他接口可用于特定用途。



#### 1.2 BeanFactory集成了环境配置

> 这种方法的重点是 `BeanFactory` 是应用程序组件的注册中心，并且它集成了应用程序组件的配置（例如不再需要单个对象读取属性文件）。有关此方法的好处的讨论，请参见《Expert One-on-One J2EE Design and Development》的第4章和第11章。

这部分解释了 `BeanFactory` 它本身是所有 Bean 的注册中心，所有的 Bean 最终都在 `BeanFactory` 中创建和保存。另外 `BeanFactory` 中还集成了配置信息，这部分咱在第 8 章依赖注入中有接触到，咱通过加载外部的 properties 文件，借助 SpringFramework 的方式将配置文件的属性值设置到 Bean 对象中。

spring 3.1 之后就有个新的概念`Environment` , 它真正做环境和配置保存地



#### 1.3 BeanFactory推荐使用DI而不是DL

> 请注意，通常最好使用依赖注入（“推”的配置），通过setter方法或构造器注入的方式，配置应用程序对象，而不是使用任何形式的“拉”的配置（例如借助 `BeanFactory` 进行依赖查找）。 SpringFramework 的 Dependency Injection 功能是使用 `BeanFactory` 接口及其子接口实现的。

SpringFramework 官方在 IOC 的两种实现上的权衡：**推荐使用 DI ，尽可能不要使用 DL** 。

另外它这里面的一个概念特别好：**DI 的思想是“推”**，它主张把组件需要的依赖“推”到组件的成员上；**DL 的思想是”拉“**，组件需要哪些依赖需要组件自己去 IOC 容器中“拉取”。这样在解释 DL 和 DI 的概念和对比时就有了新的说法；



#### 1.4 BeanFactory支持多种类型的配置源

> 通常情况下，`BeanFactory` 会加载存储在配置源（例如 XML 文档）中 bean 的定义，并使用 `org.springframework.beans` 包中的 API 来配置 bean 。然而，`BeanFactory` 的实现可以根据需要直接在 Java 代码中返回它创建的 Java 对象。bean 定义的**存储方式没有任何限制**，它可以是 LDAP （轻型文件目录访问协议），RDBMS（关系型数据库系统），XML，properties 文件等。鼓励实现以支持 Bean 之间的引用（依赖注入）。

这一段告诉我们，SpringFramework 可以支持的配置源类型有很多种，当然咱最常用的还是 xml 和注解驱动啦 ~ 这些配置源中存储的信息是一些 Bean 的定义，



#### 1.5 BeanFactory可实现层次性

> 与 `ListableBeanFactory` 中的方法相比，`BeanFactory` 中的所有操作还将检查父工厂（如果这是 `HierarchicalBeanFactory` ）。如果在 `BeanFactory` 实例中没有找到指定的 bean ，则会**向父工厂中搜索查找**。`BeanFactory` 实例中的 Bean 应该**覆盖任何父工厂中的同名 Bean 。**

这部分想告诉我们的是，`BeanFactory` 本身可以支持**父子结构**，这个父子结构的概念和实现由 `HierarchicalBeanFactory` 实现，在 `BeanFactory` 中它也只是提了一下。

#### 1.6 BeanFactory中设有完整的生命周期控制机制

> `BeanFactory` 接口实现了尽可能支持标准 Bean 的生命周期接口。全套初始化方法及其标准顺序为：......
>
> 在关闭 `BeanFactory` 时，以下生命周期方法适用：......

Bean 的生命周期是在 `BeanFactory` 中就有设计的，而且官方文档也提供了全套的初始化和销毁流程

#### 1.7 作用域

> `BeanFactory` 接口由包含多个 bean 定义的对象实现，每个 bean 的**定义信息均由 “name” 进行唯一标识**。根据 bean 的定义，SpringFramework 中的工厂会返回所**包含对象的独立实例** ( prototype ，原型模式 ) ，或者**返回单个共享实例** ( singleton ，单例模式的替代方案，其中实例是工厂作用域中的单例 ) 。返回 bean 的实例类型取决于 bean 工厂的配置：API是相同的。从 SpringFramework 2.0 开始，根据具体的应用程序上下文 ( 例如 Web 环境中的 request 和 session 作用域 ) ，可以使用更多作用域。

默认情况下，`BeanFactory` 中的 Bean 只有**单实例 Bean（`singleton`）** 和**原型 Bean（`prototype`）** ，自打 SpringFramework2.0 开始，出现了 Web 系列的作用域 `“request”` 和 `“session”` ，后续的又出现了 `“global session”` 和 `“websocket”` 作用域。

解读：

> 返回 bean 的实例类型取决于 bean 工厂的配置：API是相同的。

```java
@Component
@Scope("prototype")
public class Cat { }
```

无论是声明单实例 Bean ，还是原型 Bean ，都是用 `@Scope` 注解标注；在配置类中用 `@Bean` 注册组件，如果要显式声明作用域，也是用 `@Scope` 注解。由此就可以解释这句话了：**产生单实例 Bean 和原型 Bean 所用的 API 是相同的，都是用 `@Scope` 注解来声明，然后由 `BeanFactory` 来创建**。

#### 1.8 小结

`BeanFactory` 提供了如下基础的特性：

- 基础的容器
- 定义了作用域的概念
- 集成环境配置
- 支持多种类型的配置源
- 层次性的设计
- 完整的生命周期控制机制



### 2. HierarchicalBeanFactory

从类名上能很容易的理解，它是体现了**层次性**的 `BeanFactory` 。有了这个特性，`BeanFactory` 就有了**父子结构**。它的文档注释蛮简单的，咱看一眼：

> 由 `BeanFactory` 实现的子接口，它可以理解为是层次结构的一部分。
>
> 可以在 `ConfigurableBeanFactory` 接口中找到用于 `BeanFactory` 的相应 `setParentBeanFactory` 方法，该方法允许以可配置的方式设置父对象。

~~~java
public interface HierarchicalBeanFactory extends BeanFactory {

	/**
	 * Return the parent bean factory, or {@code null} if there is none.
	 */
	@Nullable
	BeanFactory getParentBeanFactory();

	/**
	 * Return whether the local bean factory contains a bean of the given name,
	 * ignoring beans defined in ancestor contexts.
	 * <p>This is an alternative to {@code containsBean}, ignoring a bean
	 * of the given name from an ancestor bean factory.
	 * @param name the name of the bean to query
	 * @return whether a bean with the given name is defined in the local factory
	 * @see BeanFactory#containsBean
	 */
	boolean containsLocalBean(String name);

}
~~~

**`getParentBeanFactory()`** ，它就可以获取到父 `BeanFactory` 对象；接口中还有一个方法是 `containsLocalBean(String name)` ，它是检查当前本地的容器中是否有指定名称的 Bean ，而不会往上找父 `BeanFactory` 。

`ConfigurableBeanFactory` 接口里面有`setParentBeanFactory` 设置父对象



问题：如果当前 `BeanFactory` 中有指定的 Bean 了，父 `BeanFactory` 中可能有吗？

答案是有，因为**即便存在父子关系，但他们本质上是不同的容器，所以有可能找到多个相同的 Bean** 。换句话说，**`@Scope` 中声明的 Singleton 只是在一个容器中是单实例的，但有了层次性结构后，对于整体的多个容器来看，就不是单实例的了**。



### 3.ListableBeanFactory

**可以列举出容器中的所有Bean**

> 它是 `BeanFactory` 接口的扩展实现，它可以列举出所有 bean 实例，而不是按客户端调用的要求，按照名称一一进行 bean 的依赖查找。具有 “预加载其所有 bean 定义信息” 的 `BeanFactory` 实现（例如基于XML的 `BeanFactory` ）可以实现此接口。



#### 3.1 只会列举当前容器中的Bean

> 如果当前 `BeanFactory` 同时也是 `HierarchicalBeanFactory` ，则返回值会忽略 `BeanFactory` 的层次结构，仅仅与当前 `BeanFactory` 中定义的 bean 有关。除此之外，也可以使用 `BeanFactoryUtils` 来考虑父 `BeanFactory` 中的 bean 。

如果真的想获取所有 Bean ，可以借助 `BeanFactoryUtils` 工具类来实现（工具类中有不少以 `"IncludingAncestors"` 结尾的方法，代表可以一起取父容器）。

类似`beansOfTypeIncludingAncestors` ,返回给定类型或子类型的bean



#### 3.2 有选择性的列举

> `ListableBeanFactory` 中的方法将仅遵循当前工厂的 bean 定义，它们将忽略通过其他方式（例如 `ConfigurableBeanFactory` 的 `registerSingleton` 方法）注册的**任何单实例 bean** （但 `getBeanNamesForType` 和 `getBeansOfType` 除外），它们也会检查这种手动注册的单实例 Bean 。当然，`BeanFactory` 的 `getBean` 确实也允许透明访问此类特殊 bean 。在一般情况下，无论如何所有的 bean 都来自由外部的 bean 定义信息，因此大多数应用程序不必担心这种区别。

按道理应该会把当前容器中的所有Bean都列出，但由上得**有些Bean会被忽略不会被列举出**，下面会进行解释。



##### 3.2.1 创建Bean + 配置文件

~~~java
public class Cat { }
public class Dog { }
~~~

注册Cat

~~~xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd">

<bean class="org.example.ListableBean.Cat"/>
</beans>
~~~

##### 3.2.2 驱动原始的BeanFactory加载配置文件

~~~java
public class TestMAIN {

public static void main(String[] args) {
    ClassPathResource resource = new ClassPathResource("ListableBean.xml");
    DefaultListableBeanFactory defaultListableBeanFactory = new DefaultListableBeanFactory();
    XmlBeanDefinitionReader xmlBeanDefinitionReader = new XmlBeanDefinitionReader(defaultListableBeanFactory);
    xmlBeanDefinitionReader.loadBeanDefinitions(resource);
    System.out.println("加载xml后的bean");
    Stream.of(defaultListableBeanFactory.getBeanDefinitionNames()).forEach(System.out::println);
}
}

~~~

打印出只有一个cat

~~~sh
加载xml后的bean
org.example.ListableBean.Cat#0
~~~



##### 3.2.3 手动注册Bean

~~~java
public class TestMAIN {

    public static void main(String[] args) {
        ClassPathResource resource = new ClassPathResource("ListableBean.xml");
        DefaultListableBeanFactory defaultListableBeanFactory = new DefaultListableBeanFactory();
        XmlBeanDefinitionReader xmlBeanDefinitionReader = new XmlBeanDefinitionReader(defaultListableBeanFactory);
        xmlBeanDefinitionReader.loadBeanDefinitions(resource);
        System.out.println("加载xml后的bean");
        Stream.of(defaultListableBeanFactory.getBeanDefinitionNames()).forEach(System.out::println);
        System.out.println();

        defaultListableBeanFactory.registerSingleton("dogg",new Dog());
        System.out.println("手动注册单实例后容器中的Bean");
        Stream.of(defaultListableBeanFactory.getBeanDefinitionNames()).forEach(System.out::println);

    }
}
~~~

打印得

~~~sh
加载xml后的bean
org.example.ListableBean.Cat#0

手动注册单实例后容器中的Bean
org.example.ListableBean.Cat#0
~~~

可以发现并没有打印出手动注册的`dogg`,也就是文档里面写的**忽略其他方式**；

##### 3.2.4 查看手动注册的Bean

```java
System.out.println("容器有注册Dog" + defaultListableBeanFactory.getBean("dogg"));
System.out.println("容器中所有的Dog"+ Arrays.toString(defaultListableBeanFactory.getBeanNamesForType(Dog.class)));
```

打印得


~~~sh
容器有注册Dogorg.example.ListableBean.Dog@159f197
容器中所有的Dog[dogg]
~~~

可以看到也可以取到手动注册的Bean

##### 3.2.5 设计选择性的目的

发现在一个叫 `AbstractApplicationContext` 的 `prepareBeanFactory()（容器初始化时条用，提前进行一些配置设置）` 方法中有一些使用`registerSingleton`

~~~java
protected void prepareBeanFactory(ConfigurableListableBeanFactory beanFactory) {
    ....
 		if (!beanFactory.containsLocalBean(ENVIRONMENT_BEAN_NAME)) {
			beanFactory.registerSingleton(ENVIRONMENT_BEAN_NAME, getEnvironment());
		}
		if (!beanFactory.containsLocalBean(SYSTEM_PROPERTIES_BEAN_NAME)) {
			beanFactory.registerSingleton(SYSTEM_PROPERTIES_BEAN_NAME, getEnvironment().getSystemProperties());
		}
		if (!beanFactory.containsLocalBean(SYSTEM_ENVIRONMENT_BEAN_NAME)) {
			beanFactory.registerSingleton(SYSTEM_ENVIRONMENT_BEAN_NAME, getEnvironment().getSystemEnvironment());
		}
}
~~~

> 这段代码是用于注册三个单例bean到Spring应用程序上下文中的。这些bean分别是：
>
> 1. `ENVIRONMENT_BEAN_NAME` - 用于获取应用程序的环境变量和属性配置的对象。
> 2. `SYSTEM_PROPERTIES_BEAN_NAME` - 用于获取系统属性的对象。
> 3. `SYSTEM_ENVIRONMENT_BEAN_NAME` - 用于获取系统环境变量的对象。
>
> 在这段代码中，它首先检查容器中是否已经存在这些bean，如果不存在就通过调用相应的方法创建它们，并将其注册为单例bean。这样做的目的是避免重复创建这些对象，从而提高应用程序的性能和效率。

这些组件包含一些系统变量等，spring内部使用，所以做出隐藏一些bean的功能，**目的是不希望开发者进行使用，以免对机器做出一些错误操作**

##### 3.2.6 ListableBeanFactory的大部分方法不适合频繁调用

> 注意：除了 `getBeanDefinitionCount` 和 `containsBeanDefinition` 之外，此接口中的方法不适用于频繁调用，方法的实现可能执行速度会很慢。

这个咱也能理解，毕竟谁会动不动去翻 IOC 容器的东西呢？顶多是读完一遍就自己缓存起来吧！而且一般情况下也不会有业务需求会深入到 IOC 容器的底部吧



### 4.ConfigurableBeanFactory

从类名就可以知道这是一个**可配置的**BeanFactory



#### 4.1 可读&可写

回想一开始学习面向对象编程时，就知道一个类的属性设置为 private 后，提供 **get** 方法则意味着该属性**可读**，提供 **set** 方法则意味着该属性**可写**。同样的，在 SpringFramework 的这些 `BeanFactory` ，包括后面的 `ApplicationContext` 中，都会有这样的设计。普通的 `BeanFactory` 只有 get 相关的操作，而 **Configurable** 开头的 `BeanFactory` 或者 `ApplicationContext` 就具有了 set 的操作：（节选自 `ConfigurableBeanFactory` 的方法列表）

![image-20230418103251978](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230418103251978.png)



![image-20230418103316164](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230418103316164.png)



#### 4.2 提供可配置的功能

> BeanFactory interface. 大多数 `BeanFactory` 的实现类都会实现这个带配置的接口。除了 `BeanFactory` 接口中的基础获取方法之外，还提供了配置 `BeanFactory` 的功能。



#### 4.3 不推荐给开发者使用

> `ConfigurableBeanFactory` 接口并不希望开发者在应用程序代码中使用，而是坚持使用 `BeanFactory` 或 `ListableBeanFactory` 。此扩展接口**仅用于允许在框架内部进行即插即用**，并允许对 `BeanFactory` 中的配置方法的特殊访问。

SpringFramework 不希望开发者用 `ConfigurableBeanFactory` ，而是老么实的用最根本的 `BeanFactory` ，原因也很简单，**程序在运行期间按理不应该对 `BeanFactory` 再进行频繁的变动**，此时只应该有读的动作，而不应该出现写的动作。

### 5.实现类

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/07b6807731d4430185d5ec7ea8483fbc~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



#### 5.1 AbstractBeanFactory

##### 5.1.1 AbstractBeanFactory是最终BeanFactory的基础实现

> 它是 `BeanFactory` 接口最基础的抽象实现类，提供 `ConfigurableBeanFactory` SPI 的全部功能。我们不假定有一个可迭代的 `BeanFactory` ，因此也可以用作 `BeanFactory` 实现的父类，该实现可以从某些后端资源（其中 bean 定义访问是一项昂贵的操作）获取 bean 的定义。

AbstractBeanFactory 是作为 BeanFactory 接口下面的第一个抽象的实现类，它具有最基础的功能，并且它可以从配置源（之前看到的 xml 、LDAP 、RDBMS 等）获取 Bean 的定义信息，而这个 Bean 的定义信息就是 BeanDefinition

**SPI** 全称为 **Service Provider Interface**，是 jdk 内置的一种服务提供发现机制。说白了，它可以加载预先在特定位置下配置的一些类

##### 5.1.2 AbstractBeanFactory对Bean的支持

> 此类可以提供单实例 Bean 的缓存（通过其父类 `DefaultSingletonBeanRegistry` ），单例/原型 Bean 的决定，`FactoryBean` 处理，Bean 的别名，用于子 bean 定义的 bean 定义合并以及 bean 销毁（ `DisposableBean` 接口，自定义 `destroy` 方法）。此外，它可以通过实现 `HierarchicalBeanFactory` 接口来管理 `BeanFactory` 层次结构（在未知 bean 的情况下委托给父工厂）。

根据上方接口得提供了很多功能

##### 5.1.3 AbstractBeanFactory定义了模板方法

> 子类要实现的主要模板方法是 `getBeanDefinition` 和 `createBean` ，分别为给定的 bean 名称检索 bean 定义信息，并根据给定的 bean 定义信息创建 bean 的实例。这些操作的默认实现可以在 `DefaultListableBeanFactory` 和 `AbstractAutowireCapableBeanFactory` 中找到。

SpringFramework 中大量使用**模板方法模式**来设计核心组件，它的思路是：**父类提供逻辑规范，子类提供具体步骤的实现**。在文档注释中，咱看到 `AbstractBeanFactory` 中对 `getBeanDefinition` 和 `createBean` 两个方法进行了规范上的定义，分别代表获取 Bean 的定义信息，以及创建 Bean 的实例，这两个方法都会在 SpringFramework 的 IOC 容器初始化阶段起到至关重要的作用。

~~~java
protected abstract BeanDefinition getBeanDefinition(String beanName) throws BeansException;

/**
 * Create a bean instance for the given merged bean definition (and arguments).
 * The bean definition will already have been merged with the parent definition
 * in case of a child definition.
 * <p>All bean retrieval methods delegate to this method for actual bean creation.
 * @param beanName the name of the bean
 * @param mbd the merged bean definition for the bean
 * @param args explicit arguments to use for constructor or factory method invocation
 * @return a new instance of the bean
 * @throws BeanCreationException if the bean could not be created
 */
protected abstract Object createBean(String beanName, RootBeanDefinition mbd, @Nullable Object[] args)
        throws BeanCreationException;
~~~



多说一句，`createBean` 是 SpringFramework 能管控的所有 Bean 的创建入口。



#### 5.2 AbstractAutowireCapableBeanFactory

根据类名得实现了组件的自动装配。

~~~java
public abstract class AbstractAutowireCapableBeanFactory extends AbstractBeanFactory{
    
}
~~~

##### 5.2.1 提供Bean的创建

> 它是实现了默认 bean 创建逻辑的的抽象的 `BeanFactory` 实现类，它具有 `RootBeanDefinition` 类指定的全部功能。除了 `AbstractBeanFactory` 的 `createBean` 方法之外，还实现 `AutowireCapableBeanFactory` 接口。

~~~java
protected Object createBean(String beanName, RootBeanDefinition mbd, @Nullable Object[] args)
        throws BeanCreationException {

    if (logger.isTraceEnabled()) {
        logger.trace("Creating instance of bean '" + beanName + "'");
    }
    RootBeanDefinition mbdToUse = mbd;
    。。。。
}
~~~

`AbstractAutowireCapableBeanFactory` 继承了 `AbstractBeanFactory` 抽象类，还额外实现了 `AutowireCapableBeanFactory` 接口，那实现了这个接口就代表着，它可以**实现自动注入的功能**了。除此之外，它还把 `AbstractBeanFactory` 的 `createBean` 方法给实现了，代表它还具有**创建 Bean 的功能**。



这个地方要多说一嘴，其实 **`createBean` 方法也不是最终实现 Bean 的创建**，而是有另外一个叫 **`doCreateBean`** 方法，它同样在 `AbstractAutowireCapableBeanFactory` 中定义，而且是 **protected** 方法，没有子类重写它，算是它独享的了。

~~~java
protected Object doCreateBean(String beanName, RootBeanDefinition mbd, @Nullable Object[] args)
        throws BeanCreationException {
~~~



##### 5.2.2 实现了属性赋值和组件注入

> 提供 Bean 的创建（具有构造方法的解析），属性填充，属性注入（包括自动装配）和初始化。处理运行时 Bean 的引用，解析托管集合，调用初始化方法等。支持自动装配构造函数，按名称的属性和按类型的属性。

这一段已经把 `AbstractAutowireCapableBeanFactory` 中实现的最最核心功能全部列出来了：**Bean 的创建、属性填充和依赖的自动注入、Bean 的初始化**。这部分是**创建 Bean 最核心的三个步骤**



##### 5.2.3 保留了模板方法

> 子类要实现的主要模板方法是 `resolveDependency(DependencyDescriptor, String, Set, TypeConverter)` ，用于按类型自动装配。如果工厂能够搜索其 bean 定义，则通常将通过此类搜索来实现匹配的 bean 。对于其他工厂样式，可以实现简化的匹配算法。

跟 `AbstractBeanFactory` 不太一样，`AbstractAutowireCapableBeanFactory` 没有把全部模板方法都实现完，它保留了文档注释中提到的 `resolveDependency` 方法，这个方法的作用是**解析 Bean 的成员中定义的属性依赖关系**。

~~~java
@Override
@Nullable
public Object resolveDependency(DependencyDescriptor descriptor, @Nullable String requestingBeanName) throws BeansException {
    return resolveDependency(descriptor, requestingBeanName, null, null);
}
~~~



##### 5.2.4  不负责BeanDefinition的注册

`AbstractAutowireCapableBeanFactory` 实现了对 Bean 的创建、赋值、注入、初始化的逻辑，但对于 Bean 的定义是如何进入 `BeanFactory` 的，它不负责。这里面涉及到两个流程：**Bean 的创建**、**Bean 定义的进入**，这个咱放到后面 `BeanDefinition` 和 Bean 的完整生命周期中再详细解释。



#### 5.3 DefaultListableBeanFactory

这个类是**唯一一个目前使用的 `BeanFactory` 的落地实现了**

> Spring 的 `ConfigurableListableBeanFactory` 和 `BeanDefinitionRegistry` 接口的默认实现，它时基于 Bean 的定义信息的的成熟的 `BeanFactory` 实现，它可通过后置处理器进行扩展。

 翻看源码就知道，`DefaultListableBeanFactory` 已经没有 **abstract** 标注了，说明它可以算作一个**成熟的落地实现**了。



##### 5.3.1 会先注册Bean定义信息再创建Bean

> 典型的用法是在访问 bean 之前**先注册所有 bean 定义信息**（可能是从有 bean 定义的文件中读取）。因此，按名称查找 Bean 是对本地 Bean 定义表进行的合理操作，该操作对预先解析的 Bean 定义元数据对象进行操作。

`DefaultListableBeanFactory` 在 `AbstractAutowireCapableBeanFactory` 的基础上，完成了**注册 Bean 定义信息**的动作，而这个动作就是通过上面的 **`BeanDefinitionRegistry`** 来实现的。

~~~java
void registerBeanDefinition(String beanName, BeanDefinition beanDefinition)
        throws BeanDefinitionStoreException;

/**
 * Remove the BeanDefinition for the given name.
 * @param beanName the name of the bean instance to register
 * @throws NoSuchBeanDefinitionException if there is no such bean definition
 */
void removeBeanDefinition(String beanName) throws NoSuchBeanDefinitionException;
~~~

对Bean的管理流程，**先注册 Bean 的定义信息，再完成 Bean 的创建和初始化动作**。

##### 5.3.2 不负责解析Bean定义文件

> 请注意，特定 bean 定义信息格式的解析器通常是单独实现的，而不是作为 `BeanFactory` 的子类实现的，有关这部分的内容参见 `PropertiesBeanDefinitionReader` 和 `XmlBeanDefinitionReader` 。

`BeanFactory` 作为一个统一管理 Bean 组件的容器，它的核心工作就是**控制 Bean 在创建阶段的生命周期**，而对于 Bean 从哪里来，如何被创建，都有哪些依赖要被注入，这些统统与它无关，而是有专门的组件来处理（就是包括上面提到的 `BeanDefinitionReader` 在内的一些其它组件）。

**对xml进行解析**

~~~java
public class XmlBeanFactory extends DefaultListableBeanFactory {

private final XmlBeanDefinitionReader reader = new XmlBeanDefinitionReader(this);
~~~



##### 5.3.3 替代实现

`StaticListableBeanFactory` ，它实现起来相对简单且功能也简单，因为它只能管理单实例 Bean ，而且没有跟 Bean 定义等相关的高级概念在里面，于是 SpringFramework 默认也不用它。



#### 5.4 XmlBeanFactory

在 SpringFramework 3.1 之后，`XmlBeanFactory` 正式被标注为**过时**，代替的方案是使用 `DefaultListableBeanFactory + XmlBeanDefinitionReader` ，这种设计更**符合组件的单一职责原则**

前面的demo

~~~java
ClassPathResource resource = new ClassPathResource("ListableBean.xml");
DefaultListableBeanFactory defaultListableBeanFactory = new DefaultListableBeanFactory();
XmlBeanDefinitionReader xmlBeanDefinitionReader = new XmlBeanDefinitionReader(defaultListableBeanFactory);
~~~



自打 SpringFramework 3.0 之后出现了注解驱动的 IOC 容器，SpringFramework 就感觉这种 xml 驱动的方式不应该单独成为一种方案了，倒不如咱都各退一步，**搞一个通用的容器，都组合它来用**，这样就实现了**配置源载体分离**的目的了。



## 8.IOC容器的详细对比-ApplicationContext

推荐使用 `ApplicationContext` 而不是 `BeanFactory` ，因为 `ApplicationContext` 相比较 `BeanFactory` 扩展的实在是太多了：

![image-20230418160245860](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230418160245860.png)



![img](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b31a2ee5069f4234abd048c893126ed0~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



### 8.1 ApplicationContext

#### 8.1.1 最核心的接口

> 它是为应用程序提供配置的中央接口。在应用程序运行时，它是只读的，但是如果受支持的话，它可以**重新加载**。

#### 8.1.2 组合多个接口的功能

> `ApplicationContext` 提供：
>
> - 用于访问应用程序组件的 Bean 工厂方法。继承自 `ListableBeanFactory` 。
> - 以通用方式加载文件资源的能力。继承自 `ResourceLoader` 接口。
> - 能够将事件发布给注册的监听器。继承自 `ApplicationEventPublisher` 接口。
> - 解析消息的能力，支持国际化。继承自 `MessageSource` 接口。
> - 从父上下文继承。在子容器中的定义将始终优先。例如，这意味着整个 Web 应用程序都可以使用单个父上下文，而每个 servlet 都有其自己的子上下文，**该子上下文独立于任何其他 servlet 的子上下文**。

`ApplicationContext` 也是支持层级结构的，但这里它的描述是**父子上下文**，这个概念要区分理解。**上下文中包含容器，但又不仅仅是容器。容器只负责管理 Bean ，但上下文中还包括动态增强、资源加载、事件监听机制等多方面扩展功能。**

#### 8.1.3 负责部分回调注入

> 除了标准的 `BeanFactory` 生命周期功能外，`ApplicationContext` 实现还检测并调用 `ApplicationContextAware` bean 以及 `ResourceLoaderAware` bean， `ApplicationEventPublisherAware` 和 `MessageSourceAware` bean。

- `ResourceLoader` → `ResourceLoaderAware`
- `ApplicationEventPublisher` → `ApplicationEventPublisherAware`
- `MessageSource` → `MessageSourceAware`

是不是突然明白了什么？这些 Aware 注入的最终结果还是 **`ApplicationContext`** 本身啊！



### 8.2 ConfigurableApplicationContext

它也给 `ApplicationContext` 提供了 **“可写”** 的功能

#### 8.2.1 提供了可配置的可能

> 它是一个支持 SPI 的接口，它会被大多数（如果不是全部）应用程序上下文的落地实现。除了 `ApplicationContext` 接口中的应用程序上下文客户端方法外，还提供了用于配置应用程序上下文的功能。

`ConfigurableApplicationContext` 给 `ApplicationContext` 添加了用于配置的功能，这个说法可以从接口方法中得以体现。`ConfigurableApplicationContext` 中扩展了 `setParent` 、`setEnvironment` 、`addBeanFactoryPostProcessor` 、`addApplicationListener` 等方法，都是可以改变 `ApplicationContext` 本身的方法。



#### 8.2.2 只希望被调用启动和关闭

> 配置和与生命周期相关的方法都封装在这里，以避免暴露给 `ApplicationContext` 的调用者。本接口的方法仅应由启动和关闭代码使用。

`ConfigurableApplicationContext` 本身扩展了一些方法，但是它一般情况下不希望让咱开发者调用，而是只调用启动（refresh）和关闭（close）方法。注意这个一般情况是在程序运行期间的业务代码中，但如果是为了定制化 `ApplicationContext` 或者对其进行扩展，`ConfigurableApplicationContext` 的扩展则会成为切入的主目标。



### 8.3 实现接口 EnvironmentCapable

**capable** 本意为“有能力的”，在这里解释为 **“携带/组合”** 更为合适。

**在 SpringFramework 中，以 Capable 结尾的接口，通常意味着可以通过这个接口的某个特定的方法（通常是 `getXXX()` ）拿到特定的组件。**

~~~java
public interface EnvironmentCapable {

	/**
	 * Return the {@link Environment} associated with this component.
	 */
	Environment getEnvironment();

}
~~~



#### 8.3.1 Application 具有EnvironmentCapable 的功能

`Environment` 是 SpringFramework 中抽象出来的类似于**运行环境**的**独立抽象**，它内部存放着应用程序运行的一些配置。

现阶段小伙伴可以这么理解：基于 SpringFramework 的工程，在运行时包含两部分：**应用程序本身、应用程序的运行时环境**。



#### 8.3.2 ConfigurableApplicationContext获

#### 得ConfigurableEnvironment

~~~java
@Override
ConfigurableEnvironment getEnvironment();
~~~



### 8.4 MessageSource

支持国际化的组件



### 8.5 ApplicationEventPublisher

类名可以理解为，它是**事件的发布器**。SpringFramework 内部支持很强大的事件监听机制，而 ApplicationContext 作为容器的最顶级，自然也要实现观察者模式中**广播器**的角色。文档注释中对于它的描述也是异常的简单：

> Interface that encapsulates event publication functionality. Serves as a super-interface for ApplicationContext.
>
> 封装事件发布功能的接口，它作为 `ApplicationContext` 的父接口。

### 8.6 ResourcePatternResolver

从类名理解可以解释为“**资源模式解析器**”，实际上它是**根据特定的路径去解析资源文件**的

##### 8.6.1 实现方式

`PathMatchingResourcePatternResolver`

~~~java
public class TestMAIN {
    public static void main(String[] args) throws IOException {
        PathMatchingResourcePatternResolver patternResolver = new PathMatchingResourcePatternResolver();
        Resource resource = patternResolver.getResource("application.xml");
        InputStream inputStream = resource.getInputStream();
        //Java 7 开始，可以使用 try-with-resources ,自动关注流
        try(BufferedReader bufferedReader = new BufferedReader(new InputStreamReader(inputStream))) {
            String line;
            while((line= bufferedReader.readLine())!=null)
            {
                System.out.println(line);
            }
        }
    }
}
~~~

**支持Ant路径模式匹配**

可以与任何类型的位置模式一起使用（例如 `"/WEB-INF/*-context.xml"` ）：输入模式必须与策略实现相匹配。该接口仅指定转换方法，而不是特定的模式格式。

- `/WEB-INF/*.xml` ：匹配 `/WEB-INF` 目录下的任意 xml 文件
- `/WEB-INF/**/beans-*.xml` ：匹配 `/WEB-INF` 下面任意层级目录的 `beans-` 开头的 xml 文件
- `/**/*.xml` ：匹配任意 xml 文件

**匹配类路径下的文件**

classpath*:/beans.xml。`classpath*:` 前缀告诉 Spring 在所有类路径位置中搜索文件。

## 9.ApplicationContext的实现类

![img](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/001fe910c0344849aebe654534b49cca~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



### 9.1 AbstractApplicationContext

这个类是 `ApplicationContext` **最最最最核心的实现类，没有之一**

`AbstractApplicationContext` 中定义和实现了**绝大部分应用上下文的特性和功能**

#### 9.1.1 只构建功能抽象

`AbstractApplicationContext` 的抽象实现主要是规范功能（借助模板方法），实际的动作它不管，**让子类自行去实现**。

#### 9.1.2 处理特殊类型的Bean

`ApplicationContext` 比 `BeanFactory` 强大的地方是支持更多的机制，这里面就包括了**后置处理器、监听器**等，而这些器，说白了也都是**一个一个的 Bean** ，`BeanFactory` 不会把它们区别对待，但是 `ApplicationContext` 就可以区分出来，并且赋予他们发挥特殊能力的机会。

#### 9.1.3 转换多种类型

> AbstractApplicationContext可以转换为多种类型的应用上下文，这意味着它可以被子类化和扩展来支持各种不同的场景和需求。具体而言，它可以被转换为以下几种类型：
>
> 1. ClassPathXmlApplicationContext：从类路径加载XML配置文件创建应用上下文。
> 2. FileSystemXmlApplicationContext：从文件系统加载XML配置文件创建应用上下文。
> 3. AnnotationConfigApplicationContext：基于Java注解的方式创建应用上下文。
> 4. XmlWebApplicationContext：用于Web应用程序的特殊ApplicationContext，它将XML文件加载到Web应用程序的ServletContext中。
> 5. GenericApplicationContext：通用的应用上下文，适用于任何类型的应用程序。



#### 9.1.4 提供默认加载文件策略

默认情况下，`AbstractApplicationContext` 加载资源文件的策略是直接继承了 `DefaultResourceLoader` 的策略，从类路径下加载；

但在 Web 项目中，可能策略就不一样了，它可以从 `ServletContext` 中加载（扩展的子类 `ServletContextResourceLoader` 等）。

~~~java
public abstract class AbstractApplicationContext extends DefaultResourceLoader
~~~



有一个控制`ApplicationContext` 生命周期的核心方法:`refresh`

~~~java
public void refresh() throws BeansException, IllegalStateException {
    synchronized (this.startupShutdownMonitor) {
        // Prepare this context for refreshing.
        // 1. 初始化前的预处理
        prepareRefresh();

        // Tell the subclass to refresh the internal bean factory.
        // 2. 获取BeanFactory，加载所有xml配置文件中bean的定义信息（未实例化）
        ConfigurableListableBeanFactory beanFactory = obtainFreshBeanFactory();

        // Prepare the bean factory for use in this context.
        // 3. BeanFactory的预处理配置
        prepareBeanFactory(beanFactory);

        try {
            // Allows post-processing of the bean factory in context subclasses.
            // 4. 准备BeanFactory完成后进行的后置处理
            postProcessBeanFactory(beanFactory);

            // Invoke factory processors registered as beans in the context.
            // 5. 执行BeanFactory创建后的后置处理器
            invokeBeanFactoryPostProcessors(beanFactory);

            // Register bean processors that intercept bean creation.
            // 6. 注册Bean的后置处理器
            registerBeanPostProcessors(beanFactory);

            // Initialize message source for this context.
            // 7. 初始化MessageSource
            initMessageSource();

            // Initialize event multicaster for this context.
            // 8. 初始化事件派发器
            initApplicationEventMulticaster();

            // Initialize other special beans in specific context subclasses.
            // 9. 子类的多态onRefresh
            onRefresh();

            // Check for listener beans and register them.
            // 10. 注册监听器
            registerListeners();
          
            //到此为止，BeanFactory已创建完成

            // Instantiate all remaining (non-lazy-init) singletons.
            // 11. 初始化所有剩下的单例Bean
            finishBeanFactoryInitialization(beanFactory);

            // Last step: publish corresponding event.
            // 12. 完成容器的创建工作
            finishRefresh();
        } // catch ......

        finally {
            // Reset common introspection caches in Spring's core, since we
            // might not ever need metadata for singleton beans anymore...
            // 13. 清除缓存
            resetCommonCaches();
        }
    }
}
~~~



### 9.2 GenericApplicationContext

咱先从注解驱动的 IOC 容器看起，`GenericApplicationContext` 已经是一个普通的类（非抽象类）了，它里面已经具备了 `ApplicationContext` 基本的所有能力了。

**`GenericApplicationContext` 中组合了一个 `DefaultListableBeanFactory` ！！！\**由此可以得到一个非常非常重要的信息：\**`ApplicationContext` 并不是继承了 `BeanFactory` 的容器，而是组合了 `BeanFactory` ！**



#### 9.2.1 借助BeanDefinitionRegistry处理特殊Bean

典型的用法是通过 `BeanDefinitionRegistry` 接口注册各种 Bean 的定义，然后调用 `refresh()` 以使用应用程序上下文语义来初始化这些 Bean（处理 `ApplicationContextAware` ，自动检测 `BeanFactoryPostProcessors` 等）。

在底层还是调用的 `DefaultListableBeanFactory` 执行 `registerBeanDefinition` 方法，



#### 9.2.2 只能刷新一次



由于 `GenericApplicationContext` 中组合了一个 `DefaultListableBeanFactory` ，而这个 `BeanFactory` 是在 `GenericApplicationContext` 的**构造方法中就已经初始化好**了，那么初始化好的 `BeanFactory` 就**不允许在运行期间被重复刷新了**。下面是源码中的实现：

~~~java
public GenericApplicationContext() {
    // 内置的beanFactory在GenericApplicationContext创建时就已经初始化好了
    this.beanFactory = new DefaultListableBeanFactory();
}

protected final void refreshBeanFactory() throws IllegalStateException {
    if (!this.refreshed.compareAndSet(false, true)) {
        // 利用CAS，保证只能设置一次true，如果出现第二次，就抛出重复刷新异常
        throw new IllegalStateException(
                "GenericApplicationContext does not support multiple refresh attempts: just call 'refresh' once");
    }
    this.beanFactory.setSerializationId(getId());
}
~~~



#### 9.2.3 替代方案是xml

> 对于 XML Bean 定义的典型情况，只需使用 `ClassPathXmlApplicationContext` 或 `FileSystemXmlApplicationContext` ，因为它们更易于设置（但灵活性较差，因为只能将从标准的资源配置文件中读取 XML Bean 定义，而不能混合使用任意 Bean 定义的格式）。在 Web 环境中，替代方案是 `XmlWebApplicationContext` 。



### 9.3 AbstractRefreshableApplicationContext

类名直译为 “可刷新的 ApplicationContext ”，它跟上面 `GenericApplicationContext` 的最大区别之一就是它**可以被重复刷新**。

> 它是 `ApplicationContext` 接口实现的抽象父类，应该支持多次调用 `refresh()` 方法，每次都创建一个新的内部 `BeanFactory` 实例。通常（但不是必须）这样的上下文将由一组配置文件驱动，以从中加载 bean 的定义信息。而不是在构造器里面进行注入

~~~java
public AbstractRefreshableConfigApplicationContext() {
}
~~~

创建BeanFactory有个`createBeanFactory`

~~~java
protected DefaultListableBeanFactory createBeanFactory() {
    return new DefaultListableBeanFactory(getInternalParentBeanFactory());
}
~~~



#### 9.3.1 刷新的核心是加载Bean定义信息

> 子类唯一需要实现的方法是 `loadBeanDefinitions` ，它在每次刷新时都会被调用。一个具体的实现应该将 bean 的定义信息加载到给定的 `DefaultListableBeanFactory` 中，通常委托给一个或多个特定的 bean 定义读取器。 注意，`WebApplicationContexts` 有一个类似的父类。

这段话告诉我们，既然是可刷新的 `ApplicationContext` ，那它里面存放的 **Bean 定义信息应该是可以被覆盖加载的**。由于 `AbstractApplicationContext` 就已经实现了 `ConfigurableApplicationContext` 接口，容器本身可以重复刷新，那么每次刷新时就应该重新加载 Bean 的定义信息，以及初始化 Bean 实例。

~~~java
public abstract class AbstractRefreshableApplicationContext extends AbstractApplicationContext {
	....
}

public abstract class AbstractApplicationContext extends DefaultResourceLoader
		implements ConfigurableApplicationContext {
    ...
}

~~~

另外它还说，在 Web 环境下也有一个类似的父类，猜都能猜到肯定是名字里多了个 Web ：`AbstractRefreshableWebApplicationContext` ，它的特征与 `AbstractRefreshableApplicationContext` 基本一致，不重复解释。

与普通的 `ApplicationContext` 相比，`WebApplicationContext` 额外扩展的是与 Servlet 相关的部分（ request 、`ServletContext` 等），`AbstractRefreshableWebApplicationContext` 内部就组合了一个 `ServletContext` ，并且支持给 Bean 注入 `ServletContext` 、`ServletConfig` 等 Servlet 中的组件。



#### 9.3.2 最终的实现类

几个内置的最终实现类，分别是基于 xml 配置的 `ClassPathXmlApplicationContext` 和 `FileSystemXmlApplicationContext` ，以及基于注解启动的 `AnnotationConfigApplicationContext` 。这些咱已经有了解了，下面也会展开来讲。



### 9.4 AbstractRefreshableConfigApplicationContext

与上面的 `AbstractRefreshableApplicationContext` 相比较，只是多了一个 **Config** ，说明它有**扩展跟配置相关的特性**。翻看方法列表，可以看到有它自己定义的 `getConfigLocations` 方法，意为“**获取配置源路径**”，由此也就证明它确实有配置的意思了。

通篇就抽出来一句话：用于添加对指定配置位置的通用处理。由于它是基于 xml 配置的 `ApplicationContext` 的父类，所以肯定需要传入配置源路径，那这个配置的动作就封装在这个 `AbstractRefreshableConfigApplicationContext` 中了。



### 9.5 AbstractXmlApplicationContext

最终 `ClassPathXmlApplicationContext` 和 `FileSystemXmlApplicationContext` 的直接父类了。

![image-20230418230654802](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230418230654802.png)

#### 9.5.1 具备基本全部功能

由于 `AbstractXmlApplicationContext` 已经接近于最终的 xml 驱动 IOC 容器的实现了，所以它应该有基本上所有的功能。又根据子类的两种不同的配置文件加载方式，说明**加载配置文件的策略是不一样的**，所以文档注释中有说子类只需要实现 `getConfigLocations` 这样的方法就好。

对于 `AbstractXmlApplicationContext` ，还有一个非常关键的部分需要咱知道，那就是加载到配置文件后如何处理。



#### 9.5.2 有loadBeanDefinitions的实现

~~~java
@Override
protected void loadBeanDefinitions(DefaultListableBeanFactory beanFactory) throws BeansException, IOException {
    // Create a new XmlBeanDefinitionReader for the given BeanFactory.
    // 借助XmlBeanDefinitionReader解析xml配置文件
    XmlBeanDefinitionReader beanDefinitionReader = new XmlBeanDefinitionReader(beanFactory);

    // Configure the bean definition reader with this context's
    // resource loading environment.
    beanDefinitionReader.setEnvironment(this.getEnvironment());
    beanDefinitionReader.setResourceLoader(this);
    beanDefinitionReader.setEntityResolver(new ResourceEntityResolver(this));

    // Allow a subclass to provide custom initialization of the reader,
    // then proceed with actually loading the bean definitions.
    // 初始化BeanDefinitionReader，后加载BeanDefinition
    initBeanDefinitionReader(beanDefinitionReader);
    loadBeanDefinitions(beanDefinitionReader);
}
~~~

可以看到，它解析 xml 配置文件不是自己干活，是**组合了一个 `XmlBeanDefinitionReader`** ，让它去解析

~~~java
protected void loadBeanDefinitions(XmlBeanDefinitionReader reader) throws BeansException, IOException {
    Resource[] configResources = getConfigResources();
    if (configResources != null) {
        reader.loadBeanDefinitions(configResources);
    }
    String[] configLocations = getConfigLocations();
    if (configLocations != null) {
        reader.loadBeanDefinitions(configLocations);
    }
}
~~~

可以看到就是调用上面文档注释中提到的 `getConfigResources` 和 `getConfigLocations` 方法，取到配置文件的路径 / 资源类，交给 `BeanDefinitionReader` 解析。



### 9.6 ClassPathXmlApplicationContext

**落地实现**

这段话写的很明白，它支持的配置文件加载位置都是 **classpath 下取**，这种方式的一个好处是：如果工程中依赖了一些其他的 jar 包，而工程启动时需要同时传入这些 jar 包中的配置文件，那 `ClassPathXmlApplicationContext` 就可以加载它们。

**支持Ant模式声明配置文件路径**

上面 `AbstractXmlApplicationContext` 中就说了，可以重写 `getConfigLocations` 方法来调整配置文件的默认读取位置，它这里又重复了一遍。除此之外它还提到了，加载配置文件的方式可以**使用 Ant 模式匹配**（比较经典的写法当属 web.xml 中声明的 `application-*.xml` ）。



#### 9.6.1 ClassPathXmlApplicationContext解析的配置文件有先后之分

> 如果有多个配置位置，则较新的 `BeanDefinition` 会覆盖较早加载的文件中的 `BeanDefinition` ，可以利用它来通过一个额外的 XML 文件有意覆盖某些 `BeanDefinition` 。



#### 9.6.2 ApplicationContext 组合使用

> 这是一个简单的一站式便利 `ApplicationContext` 。可以考虑将 `GenericApplicationContext` 类与 `XmlBeanDefinitionReader` 结合使用，以实现更灵活的上下文配置。

`ClassPathXmlApplicationContext` 继承了 `AbstractXmlApplicationContext` ，而 `AbstractXmlApplicationContext` 实际上是内部组合了一个 `XmlBeanDefinitionReader` ，所以就可以有一种组合的使用方式：利用 `GenericApplicationContext` 或者子类 `AnnotationConfigApplicationContext` ，配合 `XmlBeanDefinitionReader` ，就可以做到注解驱动和 xml 通吃了。



### 9.7 AnnotationConfigApplicationContext

> 它本身继承了 `GenericApplicationContext` ，那自然它也只能刷新一次。

注解驱动，除了 `@Component` 及其衍生出来的几个注解，更重要的是 `@Configuration` 注解，一个被 `@Configuration` 标注的类相当于一个 xml 文件。至于下面还提到的关于 JSR-330 的东西，它没有类似于 `@Component` 的东西（它只是定义了依赖注入的标准，与组件注册无关），它只是说如果一个组件 Bean 里面有 JSR-330 的注解，那它能给解析而已。



**配置也有先后之分**

> 允许使用 `register(Class ...)` 一对一注册类，以及使用 `scan(String ...)` 进行类路径的包扫描。 如果有多个 `@Configuration` 类，则在以后的类中定义的 `@Bean` 方法将覆盖在先前的类中定义的方法。这可以通过一个额外的 `@Configuration` 类来故意**覆盖某些 `BeanDefinition`** 。

初始化的两种方式：要么注册配置类，要么直接进行包扫描。由于注解驱动开发中可能没有一个主配置类，都是一上来就一堆 `@Component` ，这个时候完全可以直接声明根扫描包，进行组件扫描。



## 10.事件机制&监听器

### 1.观察者模式

**观察者模式关注的点是某**一个对象被修改 / 做出某些反应 / 发布一个信息等，会自动通知依赖它的对象（订阅者）。

**观察者、被观察主题、订阅者**。观察者（ Observer ）需要绑定要通知的订阅者（ Subscriber ），并且要观察指定的主题（ Subject ）。

### 2.spring里面的观察者模式

体现观察者模式的特性就是事件驱动和监听器。**监听器**充当**订阅者**，监听特定的事件；**事件源**充当**被观察的主题**，用来发布事件；**IOC 容器**本身也是事件广播器，可以理解成**观察者**。

**事件源、事件、广播器、监听器**。

- **事件源：发布事件的对象**

- **事件：事件源发布的信息 / 作出的动作**

- 广播器：事件真正广播给监听器的对象

  【即`ApplicationContext`】

  - `ApplicationContext` 接口有实现 `ApplicationEventPublisher` 接口，具备**事件广播器的发布事件的能力**
  - `ApplicationEventMulticaster` 组合了所有的监听器，具备**事件广播器的广播事件的能力**

- **监听器：监听事件的对象**



### 3.具体实例

#### 3.1 监听器接口

spring里面的监听器接口

~~~java
@FunctionalInterface
public interface ApplicationListener<E extends ApplicationEvent> extends EventListener {

	/**
	 * Handle an application event.
	 * @param event the event to respond to
	 */
	void onApplicationEvent(E event);

}
~~~

自定义监听器就只需要实现这个接口里面的方法即可

**事件**

![image-20230419102632559](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230419102632559.png)

代表容器的开始，刷新开始，刷新结束，停止



**自定义一个监听器**

~~~java
@Component
public class ContextRefreshedApplicationListener implements ApplicationListener<ContextRefreshedEvent> {
    @Override
    public void onApplicationEvent(ContextRefreshedEvent event) {
        System.out.println("开始监听ContextRefreshedEvent");
    }
}
~~~

记得使用`@Component` 进行注册Bean



#### 3.2 编写启动类

~~~java
System.out.println("准备初始化IOC容器。。。");
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext("org.example.Observer");

System.out.println("IOC容器初始化完成。。。");
ctx.close();
System.out.println("IOC容器关闭。。。");
/*
准备初始化IOC容器。。。
开始监听ContextRefreshedEvent
IOC容器初始化完成。。。
IOC容器关闭。。。
*/
~~~



#### 3.3 注解式设计

~~~java
@Component
public class ContextCloseApplicationListener {

    @EventListener
    // 注意传入需要监听的事件，否则不知道监听什么而报错
    public void onClosedLister(ContextClosedEvent contextClosedEvent) {
        System.out.println("context销毁了...");
    }
}

~~~

启动类测试

~~~sh
准备初始化IOC容器。。。
#
开始监听ContextRefreshedEvent

IOC容器初始化完成。。。

#
context销毁了...
IOC容器关闭。。。
~~~



#### 3.4 小节

由这两种监听器的 Demo ，可以得出几个结论：

- `ApplicationListener` 会在容器初始化阶段就准备好，在容器销毁时一起销毁；
- `ApplicationListener` 也是 IOC 容器中的普通 Bean ；
- IOC 容器中有内置的一些事件供我们监听。



### 4.spring内置的事件



#### 4.1 ApplicationEcent

事件模型的抽象。由所有应用程序事件扩展的类。它被设计为抽象的，因为**直接发布一般事件没有意义**。

~~~java
public abstract class ApplicationEvent extends EventObject
~~~



#### 4.2 ApplicationContextEvent

~~~java
public abstract class ApplicationContextEvent extends ApplicationEvent {

	/**
	 * Create a new ContextStartedEvent.
	 * @param source the {@code ApplicationContext} that the event is raised for
	 * (must not be {@code null})
	 */
	public ApplicationContextEvent(ApplicationContext source) {
		super(source);
	}

	/**
	 * Get the {@code ApplicationContext} that the event was raised for.
	 */
	public final ApplicationContext getApplicationContext() {
		return (ApplicationContext) getSource();
	}

}
~~~

可以**通过监听器直接取到 `ApplicationContext` 而不需要做额外的操作**



#### 4.3  ContextRefreshedEvent&ContextClosedEvent

这两个是一对，分别对应着 **IOC 容器刷新完毕但尚未启动**，以及 **IOC 容器已经关闭但尚未销毁所有 Bean** 。



#### 4.4 ContextStartedEvent&ContextStoppedEvent

`ContextRefreshedEvent` 事件的触发是所有单实例 Bean 刚创建完成后，就发布的事件，此时那些实现了 `Lifecycle` 接口的 Bean 还没有被回调 `start` 方法。当这些 `start` 方法被调用后，`ContextStartedEvent` 才会被触发。同样的，`ContextStoppedEvent` 事件也是在 `ContextClosedEvent` 触发之后才会触发，**此时单实例 Bean 还没有被销毁，要先把它们都停掉才可以释放资源，销毁 Bean 。**



### 5.自定义事件开发



#### 5.1 场景描述

用户注册成功后，网站进行发送注册成功的消息。网站广播一个"用户注册成功"的事件，然后由网站的信息监听器监听到，然后通过不同方式进行传播，eg：短信，微信等



#### 5.2 自定义注册成功事件

~~~java
public class RegisterSuccessEvent extends ApplicationEvent {
    /**
     * Create a new {@code ApplicationEvent}.
     *
     * @param source the object on which the event initially occurred or with
     *               which the event is associated (never {@code null})
     */
    public RegisterSuccessEvent(Object source) {
        super(source);
    }
}
~~~



#### 5.3 自定义监听器

~~~java
@Component
public class MessageListener implements ApplicationListener {
    @Override
    public void onApplicationEvent(ApplicationEvent event) {
        System.out.println("注册成功，开始发送短信");
    }
}
~~~



~~~java
@Component
public class QQListener {

    @EventListener
    public void onSuccessListener() {
        System.out.println("发送QQ");
    }
}
~~~



#### 5.4 自定义注册逻辑

~~~java
@Service
public class RegisterService implements ApplicationEventPublisherAware {

    ApplicationEventPublisher publisher;

    public void register(String username) {
        System.out.println("注册");
        // 进行广播
        publisher.publishEvent(new RegisterSuccessEvent(username));
    }

    // 获得publisher
    @Override
    public void setApplicationEventPublisher(ApplicationEventPublisher applicationEventPublisher) {
        this.publisher = applicationEventPublisher;
    }
}

~~~



#### 5.5 启动类

~~~java
public class MAIN {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext("org.example.Observer");
        RegisterService bean = ctx.getBean(RegisterService.class);
        bean.register("zqy");
    }
}
/*
注册成功，开始发送短信
(前面的内置监听器)
开始监听ContextRefreshedEvent
注册
发送QQ
注册成功，开始发送短信
*/
~~~

#### 5.6 监听器执行顺序

调整执行顺序，可以使用`@Order`注解，默认的排序值为 **`Integer.MAX_VALUE`** ，代表**最靠后**

> 不加@Order默认是0；
> @Order默认是Integer.MAX_VALUE；
> 所以@Order(-1)比不加注解先执行



## 11.模板装配

### 1.原生手动进行装配

使用 `@Configuration` + `@Bean` 注解组合，或者 `@Component` + `@ComponentScan` 注解组合，可以实现编程式 / 声明式的手动装配。

但是如果Bean多了，那么需要添加注解的地方就会很多，导致很麻烦



### 2.模板装配

#### 2.1 模板

模块可以理解成一个一个的可以分解、组合、更换的独立的单元，模块与模块之间可能存在一定的依赖，**模块的内部通常是高内聚的**，一个模块通常都是解决一个独立的问题

- 独立的
- 功能高内聚
- 可相互依赖
- 目标明确

#### 2.2 模板装配

理解为**把一个模块需要的核心功能组件都装配好**

#### 2.3 spring内的模板装配

引入大量 **`@EnableXXX`** 注解，来快速整合激活相对应的模块

- `EnableTransactionManagement `：开启注解事务驱动
- `EnableWebMvc` ：激活 SpringWebMvc
- `EnableAspectJAutoProxy` ：开启注解 AOP 编程
- `EnableScheduling` ：开启调度功能（定时任务）



#### 2.4 实践

自定义注解+`@Import` 导入组件



##### 2.4.1 设置场景

构建出一个**酒馆**，酒馆里得有**吧台**，得有**调酒师**，得有**服务员**，还得有**老板**。

模拟实现的最终目的，是可以**通过一个注解，同时把这些元素都填充到酒馆中**。



##### 2.4.2 声明自定义注解

~~~java
@Documented
@Retention(RetentionPolicy.RUNTIME)
@Target(ElementType.TYPE)
@Import
public @interface EnableTavern {
}

~~~

`@Import`，查看源码得有一个value属性，指定要导入的Configuration,ImportSelector,ImportBeanDefinitionRegistrar,或者普通Bean

~~~java
public @interface Import {

	/**
	 * {@link Configuration @Configuration}, {@link ImportSelector},
	 * {@link ImportBeanDefinitionRegistrar}, or regular component classes to import.
	 */
	Class<?>[] value();

}
~~~



##### 2.4.3 声明Bean引入&配置类

~~~java
public class Boss {

}
===========
@Import({Boss.class})
public @interface EnableTavern {
}
~~~

**配置类**

~~~java
@Configuration
@EnableTavern
public class TavernConfiguration {

}
~~~

##### 2.4.4 启动类测试

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(TavernConfiguration.class);
String[] names = ctx.getBeanDefinitionNames();
Stream.of(names).forEach(System.out::println);

/*
tavernConfiguration
org.example.FormworkAssembly.normalAss.Boss
*/
~~~



#### 2.5 模板配置的方式

##### 2.5.1 导入普通类

如上

##### 2.5.2 导入配置6类

~~~java
public class Bartender {
    
    private String name;
    
    public Bartender(String name) {
        this.name = name;
    }
    
    public String getName() {
        return name;
    }
}
~~~

 

~~~java
@Configuration
public class BartenderConfiguration {
    
    @Bean
    public Bartender zhangxiaosan() {
        return new Bartender("张小三");
    }
    
    @Bean
    public Bartender zhangdasan() {
        return new Bartender("张大三");
    }
    
}
~~~

~~~java
@Import({Boss.class, BartenderConfiguration.class})
~~~

> 注意这里有一个小细节，有小伙伴在学习的时候，启动类里或者配置类上用了**包扫描**，恰好把这个类扫描到了，导致即使没有 `@Import` 这个 `BartenderConfiguration` ，`Bartender` 调酒师也被注册进 IOC 容器了。这里一定要细心哈，包扫描本身就会扫描配置类，并且让其生效的。如果既想用包扫描，又不想扫到这个类，很简单，把这些配置类拿到别的包里，让包扫描找不到它就好啦。



 **测试运行**

~~~sh
tavernConfiguration
org.example.FormworkAssembly.normalAss.Boss
org.example.FormworkAssembly.normalAss.BartenderConfiguration
zhangxiaosan
zhangdasan
~~~



##### 2.5.3  导入ImportSelector

> 它是一个接口，它的实现类可以根据指定的筛选标准（通常是一个或者多个注解）来决定导入哪些配置类。

~~~java
public class Bar {
    
}
~~~



**配置类**

~~~java
@Configuration
public class BarConfiguration {
    
    @Bean
    public Bar bbbar() {
        return new Bar();
    }
}
~~~



**实现ImportSelector**

~~~java
public class BarImportSelector implements ImportSelector {
    @Override
    public String[] selectImports(AnnotationMetadata importingClassMetadata) {
        return new String[0];
    }
}
~~~

这里返回的字符串数组是**全限定类名**，通过全限定类名进行定位具体的类，然后进行注册。



**注入@ImportSelector&测试**

~~~java
org.example.FormworkAssembly.normalAss.Bar
org.example.FormworkAssembly.normalAss.BarConfiguration
bbbar
~~~

可以看到`ImportSelector`并没有被注册到IOC容器



##### 2.5.4 导入ImportBeanDefinitionRegistrar

导入的实际是 `BeanDefinition` （ Bean 的定义信息）

~~~java
public class Waiter {
    
}
~~~



**编写实现类**

~~~java
public class WaiterRegistrar implements ImportBeanDefinitionRegistrar {
    @Override
    public void registerBeanDefinitions(AnnotationMetadata importingClassMetadata, BeanDefinitionRegistry registry) {
        registry.registerBeanDefinition("waiter", new RootBeanDefinition(Waiter.class));
    }
}
~~~



**启动类测试**

~~~java
//waiter
~~~

如上没有注册`WaiterRegistrar`



**综上模板匹配就是这四种方法进行组合使用**



### 3.条件装配

上面的模板装配，只要配置类里面声明了`@Bean` ,方法的返回直接就会被注册到IOC容器称为一个Bean。**如何根据自己需要进行注册Bean？**

#### 1.profile

##### 1.1 什么是profile

profile 有“配置文件”的意思，倒不是说一个 profile 是一个配置文件，它**更像是一个标识**。



`@Profile` 注解可以标注在组件上，当一个配置属性（并不是文件）激活时，它才会起作用，而激活这个属性的方式有很多种（启动参数、环境变量、`web.xml` 配置等）。

**根据当前项目运行的环境不同，可以动态的注册当前运行环境需要的组件**



##### 1.2 使用

添加`@Profile`

~~~java
@Configuration
@Profile(value = "city")
public class BartenderConfiguration {

    @Bean
    public Bartender zhangxiaosan() {
        return new Bartender("张小三");
    }

    @Bean
    public Bartender zhangdasan() {
        return new Bartender("张大三");
    }

}
~~~

启动类测试得

```sh
org.example.FormworkAssembly.normalAss.Boss
org.example.FormworkAssembly.normalAss.Bar
org.example.FormworkAssembly.normalAss.BarConfiguration
bbbar
waiter
```

配置类里面的Bartender没有注册到IOC容器。

默认情况下，`ApplicationContext` 中的 profile 为 **“default”**，那上面 `@Profile("city")` 不匹配，`BartenderConfiguration` 不会生效，



###### 1.2.1 编程式设置运行环境

~~~java
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(TavernConfiguration.class);
        ctx.getEnvironment().setActiveProfiles("city");
        Stream.of(ctx.getBeanDefinitionNames()).forEach(System.out::println);
~~~

但运行后还是发现Bartender**还是没有注册到IOC容器**

---

> 为什么没有注册？

前面学习`ApplicationContext` 里面有个`refresh` 方法，这里在`new AnnotationConfigApplicationContext()` 就将配置类，内部就已经初始化了，所以即使后面再设置环境也没有用。

> 后续再进行传入配置类呢？

~~~java
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext();
        ctx.getEnvironment().setActiveProfiles(("city"));
        ctx.register(TavernConfiguration.class);
        ctx.refresh();
        Stream.of(ctx.getBeanDefinitionNames()).forEach(System.out::println);
/*
..
org.example.FormworkAssembly.normalAss.BartenderConfiguration
zhangxiaosan
zhangdasan
org.example.FormworkAssembly.normalAss.Bar
...
*/
~~~

可以看到Bartender打印出来了



###### 1.2.2声明式设置运行环境

![image-20230420155935405](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230420155935405.png)

~~~sh
-Dspring.profiles.active=city
~~~

改回main代码

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(TavernConfiguration.class);
//        ctx.getEnvironment().setActiveProfiles("city");
Stream.of(ctx.getBeanDefinitionNames()).forEach(System.out::println);
/*
..
org.example.FormworkAssembly.normalAss.BartenderConfiguration
zhangxiaosan
zhangdasan
org.example.FormworkAssembly.normalAss.Bar
...
*/
~~~

可以看到这样的设置，也会实现功能

##### 1.3 @Profile在实际开发的用途

不同环境连接的数据库都是不一样。为解决不同环境更改配置文件，就使用`@Profile` 来进行设置不同环境的数据源。

~~~java
@Configuration
public class DataSourceConfiguration {

    @Bean
    @Profile("dev")
    public DataSource devDataSource() {
        return null;
    }

    @Bean
    @Profile("test")
    public DataSource testDataSource() {
        return null;
    }

    @Bean
    @Profile("prod")
    public DataSource prodDataSource() {
        return null;
    }
}
~~~

然后通过`@PropertySource` + 外部配置文件，就可以做到不同环境通过profile注解切换不同数据源



##### 1.4 profile的弊端

profile有些地方无法控制。

**吧台应该是由老板安置好的，如果酒馆中连老板都没有，那吧台也不应该存在。**

因为 **profile 控制的是**整个项目的运行**环境**，**无法根据单个 Bean 的因素决定是否装配**

出现了第二种条件装配：**`@Conditional` 注解**。



#### 2.Conditional

被标注 `@Conditional` 注解的 Bean 要注册到 IOC 容器时，必须全部满足 `@Conditional` 上指定的所有条件才可以。

> 介绍：`@Conditional` 注解可以指定匹配条件，而被 `@Conditional` 注解标注的 组件类 / 配置类 / 组件工厂方法 必须满足 `@Conditional` 中指定的所有条件，才会被创建 / 解析。



##### 2.1  使用

~~~java
@Configuration
public class BarConfiguration {

    @Bean
    // 这里需要设置一个value值，条件类
    @Conditional(value = ExistBossCondition.class)
    public Bar bbbar() {
        return new Bar();
    }
}
~~~

**ExistBossCondition**

~~~java
// 注意这里的包印的是spring的，而不是...lock
import org.springframework.context.annotation.Condition;

public class ExistBossCondition implements Condition {
    @Override
    public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
       
        return context.getBeanFactory().containsBeanDefinition(Boss.class.getName())
    }
}
~~~

注意上面判断的是`containsBeanDefinition`,而不是Bean，因为考虑到构建阶段，Boss可能还没有被创建，所以就判断BeanDefinition是否存在。



##### 2.2 测试

~~~java
tavernConfiguration
org.example.FormworkAssembly.ConditionalAssembly.Boss
org.example.FormworkAssembly.ConditionalAssembly.Bar
org.example.FormworkAssembly.ConditionalAssembly.BarConfiguration
bbbar
waiter
~~~

---

怎么判断`@Condition`是否生效

将`Boss.class` 删去，再运行得

~~~java
@Import({ BartenderConfiguration.class, BarImportSelector.class, WaiterRegistrar.class})
public @interface EnableTavern {
}

~~~

可以发现该注解生效了

~~~sh
tavernConfiguration
org.example.FormworkAssembly.ConditionalAssembly.Bar
org.example.FormworkAssembly.ConditionalAssembly.BarConfiguration
waiter
~~~



#### 3. 通用抽取

如果一个项目中，有比较多的组件需要依赖另一些不同的组件，如果每个组件都写一个 `Condition` 条件，那工程量真的太大了。就想办法将匹配的规则抽取为通用的方式。

##### 3.1  抽取传入的beanName

`@Conditional` 可以派生，自定义一个新的注解`@ConditionalOnBean` ,意为**存在指定Bean时匹配**

~~~java
@Documented
@Retention(RetentionPolicy.RUNTIME)
@Target({ElementType.TYPE, ElementType.METHOD})
@Conditional(OnBeanCondition.class)
public @interface ConditionalOnBean {
    // 指定bean的名称
    String[] beanNames() default {};

}

~~~

`OnBeanCondition`

~~~java
public class OnBeanCondition implements Condition {
    @Override
    public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
        String[] beanNames = (String[]) metadata.getAnnotationAttributes(ConditionalOnBean.class.getName()).get("beanNames");
        for (String beanName : beanNames) {
            if (!context.getBeanFactory().containsBeanDefinition(beanName)) {
                return false;
            }
        }
        return true;
    }
}
~~~



启动类测试得

~~~sh
tavernConfiguration
org.example.FormworkAssembly.ConditionalAssembly.Boss
org.example.FormworkAssembly.ConditionalAssembly.Bar
org.example.FormworkAssembly.ConditionalAssembly.BarConfiguration
bbbar
waiter
~~~

因为EnableTavern里面有`Boss.class` ，所以bbbar注册到了IOC容器

##### 3.2  加入类型匹配

上面只能是抽取 `beanName` ，传整个类的全限定名真的很费劲。如果当前类路径下本来就有这个类，那直接写进去就好呀。

希望的效果为

```java
@Bean
@ConditionalOnBean(Boss.class)
public Bar bbbar() {
    return new Bar();
}
```



给 `@ConditionalOnBean` 注解上添加默认的 `value` 属性，类型为 `Class[]` ，这样就可以传入类型了：

~~~java
@Documented
@Retention(RetentionPolicy.RUNTIME)
@Target({ElementType.TYPE, ElementType.METHOD})
@Conditional(OnBeanCondition.class)
public @interface ConditionalOnBean {
    String[] beanNames() default {};
	
    Class<?>[] value() default {};

}
~~~

然后修改`OnBeanCondition`

~~~java
public class OnBeanCondition implements Condition {
    @Override
    public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
        Map<String, Object> attributes = metadata.getAnnotationAttributes(ConditionalOnBean.class.getName());

        Class<?>[] value = (Class<?>[]) attributes.get("value");
        for (Class<?> aClass : value) {
            if (!context.getBeanFactory().containsBeanDefinition(aClass.getName())) {
                return false;
            }
        }
        String[] beanNames = (String[]) attributes.get("beanNames");
        for (String beanName : beanNames) {
            if (!context.getBeanFactory().containsBeanDefinition(beanName)) {
                return false;
            }
        }
        return true;
    }
}
~~~



注意：

**使用java注解时不写属性名会默认给名为"value"的属性赋值。**



## 12. 组件扫描

### 1. 包扫描

`@ComponentScan` 注解可以指定包扫描的路径（而且还可以声明不止一个），它的写法是使用 `@ComponentScan` 的 `value` / `basePackages` 属性

~~~java
@Configuration
@ComponentScan("com.linkedbear.spring.annotation.e_basepackageclass.bean")
public class BasePackageClassConfiguration {
    
}
~~~



查看源码得还可以传入类的Class字节码

~~~java
Class<?>[] basePackageClasses() default {};
~~~

它的这个 `basePackageClasses` 属性，可以传入一组 Class 进去，它代表的意思，是扫描**传入的这些 Class 所在包及子包下的所有组件**。

#### 1.1  创建Bean

![image-20230420222221118](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230420222221118.png)

~~~java
@Service
public class DemoService {
}
~~~



~~~java
@Component
public class DemoDao {
}
~~~



~~~java
@Configuration
@ComponentScan(basePackageClasses = DemoService.class)
public class BasePackageClassConfiguration {
}
~~~

#### 1.2 **测试运行得**

传入`DemoService`的字节码，也将同一包下的`DemoBean` 也注册进去了

~~~sh
....
basePackageClassConfiguration
demoDao
demoService
~~~



### 2.  包扫描的过滤

我们用包扫描拿到的组件不一定全部都需要，也或者只有一部分需要，这个时候就需要用到包扫描的过滤了。



#### 2.1 按注解过滤包含

~~~java
@Configuration
@ComponentScan(basePackageClasses = DemoService.class,
        excludeFilters = @ComponentScan.Filter(type = FilterType.ANNOTATION, value = Service.class))
public class BasePackageClassConfiguration {
}
~~~

测试得,DemoDao还是被注册进来了

```shell
basePackageClassConfiguration
demoDao
demoService
```

因为`@ComponentScan` 注解中还有一个属性：`useDefaultFilters` ，它代表的是“是否启用默认的过滤规则”。咱之前也讲过了，默认规则就是扫那些以 `@Component` 注解为基准的模式注解。

> 指示是否应启用对以 `@Component` 、`@Repository` 、`@Service` 或 `@Controller` 注解的类的自动检测。

默认为true，设置为false后，可以看到`@Component` 注解的DemoBean并没有注册进来

~~~java
boolean useDefaultFilters() default true;
~~~

也许这样理解起来更容易一些（有颜色的部分代表匹配规则）：

![img](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/4cc6502a040548d6966aa52e89d6542a~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

#### 2.2 按照注解排除

~~~java
@Configuration
@ComponentScan(basePackageClasses = DemoService.class,
        excludeFilters = @ComponentScan.Filter(type = FilterType.ANNOTATION, value = Service.class))
public class BasePackageClassConfiguration {
}
~~~

测试得,DemoDao注册成功，但是DemoService没有了

~~~sh
basePackageClassConfiguration
demoDao
~~~

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9b8246722a7844fe9b10793033120120~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



#### 2.3  按类型过滤

FilterType.ASSIGNABLE_TYPE

~~~java
@Configuration
@ComponentScan(basePackageClasses = DemoService.class,
        excludeFilters = @ComponentScan.Filter(type = FilterType.ASSIGNABLE_TYPE, value = DemoService.class))
public class BasePackageClassConfiguration {
}

/*
...
basePackageClassConfiguration
demoDao
*/
~~~

根据DemoService类型进行过滤



#### 2.4 正则表达式过滤

~~~java
@Configuration
@ComponentScan(basePackageClasses = DemoService.class,
        excludeFilters = @ComponentScan.Filter(type = FilterType.REGEX, pattern = "org.example.ComponentScanning.bean.+Dao"))
public class BasePackageClassConfiguration {
}

~~~

解释：

`"org.example.ComponentScanning.bean.+Dao"` 排除bean文件夹下，以Dao结尾的组件，如`DemoDao`

~~~sh
...
basePackageClassConfiguration
demoService
~~~



#### 2.5 自定义过滤



##### 2.5.1 TypeFilter接口

编程时自定义过滤，需要编写过滤策略，实现 `TypeFilter` 接口。这个接口只有一个 `match` 方法：

~~~java
@FunctionalInterface
public interface TypeFilter {

	boolean match(MetadataReader metadataReader, MetadataReaderFactory metadataReaderFactory)
			throws IOException;

}
~~~

这个 `match` 方法有两个参数

- metadataReader
  - 通过这个 Reader ，可以读取到**正在扫描的类的信息**（包括类的信息、类上标注的注解等）
- metadataReaderFactory
  - 借助这个 Factory ，可以获取到**其他类的 Reader** ，进而获取到那些类的信息
  - 可以这样理解：**借助 ReaderFactory 可以获取到 Reader ，借助 Reader 可以获取到指定类的信息**

##### 2.5.2 实现自定义过滤

`MetadataReader` 中有一个 `getClassMetadata` 方法，可以拿到正在扫描的类的基本信息，咱可以由此取到全限定类名，进而与咱需求中的 类做匹配：

~~~java
public class DaoTypeFilter implements TypeFilter {
    @Override
    public boolean match(MetadataReader metadataReader, MetadataReaderFactory metadataReaderFactory) throws IOException {
        ClassMetadata classMetadata = metadataReader.getClassMetadata();
        return classMetadata.getClassName().equals(DemoDao.class.getName());
    }
}
~~~



然后在配置类上进行设置
`FilterType.CUSTOM` 表明自定义的过滤

~~~java
@Configuration
@ComponentScan(basePackageClasses = DemoService.class,
        excludeFilters = {@ComponentScan.Filter(type = FilterType.CUSTOM, value = DaoTypeFilter.class)})
public class BasePackageClassConfiguration {
}
~~~



测试得

```sh
...
basePackageClassConfiguration
demoService
```



##### 2.5.3 metadata概念

讲到这里了，咱先不着急往下走，停一停，咱讲讲 **metadata** 的概念。

回想一下 JavaSE 的反射，它是不是可以根据咱写好的类，获取到类的全限定名、属性、方法等信息呀。好，咱现在就建立起这么一个概念：咱定义的类，它叫什么名，它有哪些属性，哪些方法，这些信息，统统叫做**元信息**，**元信息会描述它的目标的属性和特征**。

在 SpringFramework 中，元信息大量出现在框架的底层设计中，不只是 **metadata** ，前面咱屡次见到的 **definition** ，也是元信息的体现。后面到了 IOC 高级部分，咱会整体的学习 SpringFramework 中的元信息、元定义设计，以及 `BeanDefinition` 的全解析。

### 3. 包扫描的其他特性

#### 3.1 包扫描组合使用

spring4.3 后出现了 `@ComponentScans` 注解

~~~java
@Retention(RetentionPolicy.RUNTIME)
@Target(ElementType.TYPE)
@Documented
public @interface ComponentScans {

	ComponentScan[] value();

}
~~~

一次性组合了一堆 `@ComponentScan`



#### 3.2 包扫描的组件名称生成

默认情况下生成的 bean 的名称是类名的首字母小写形式（ Person → person ）

实现的原因为属性`nameGenerator`

~~~java
Class<? extends BeanNameGenerator> nameGenerator() default BeanNameGenerator.class;
~~~

实现类为`AnnotaionBeanNameGenerator`

![image-20230421105233000](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230421105233000.png)



##### 3.2.1 BeanNameGenerator

~~~java
public interface BeanNameGenerator {
	String generateBeanName(BeanDefinition definition, BeanDefinitionRegistry registry);
}
~~~

又出现 `BeanDefinition` 和 `BeanDefinitionRegistry` 了，可见元信息、元定义在底层真的太常见了



##### 3.2.2 AnnotationBeanNameGenerator

~~~java
public String generateBeanName(BeanDefinition definition, BeanDefinitionRegistry registry) {
    // 组件的注册方式是注解扫描的
    if (definition instanceof AnnotatedBeanDefinition) {
        // 尝试从注解中获取名称
        String beanName = determineBeanNameFromAnnotation((AnnotatedBeanDefinition) definition);
        if (StringUtils.hasText(beanName)) {
            // Explicit bean name found.
            return beanName;
        }
    }
    // Fallback: generate a unique default bean name.
    // 如果没有获取到，则创建默认的名称
    return buildDefaultBeanName(definition, registry);
}
~~~

1. 只有注解扫描的Bean才会被处理
2. 查看注解中是否有声明

> 这种声明方式就是 `@Component("person")`

3. 查找不到就默认规则声明

##### 3.2.3 buildDefaultBeanName

~~~java
	protected String buildDefaultBeanName(BeanDefinition definition, BeanDefinitionRegistry registry) {
		return buildDefaultBeanName(definition);
	}


	protected String buildDefaultBeanName(BeanDefinition definition) {
		String beanClassName = definition.getBeanClassName();
		Assert.state(beanClassName != null, "No bean class name set");
		String shortClassName = ClassUtils.getShortName(beanClassName);
		return Introspector.decapitalize(shortClassName);
	}
~~~

用`getShortName` 截取短类名(`com.linkedbear.Person` → `Person`)

最后用一个叫 `Introspector` 的类，去生成 bean 的名称

~~~java
    public static String decapitalize(String name) {
        if (name == null || name.length() == 0) {
            return name;
        }
        if (name.length() > 1 && Character.isUpperCase(name.charAt(1)) &&
                        Character.isUpperCase(name.charAt(0))){
            return name;
        }
        char chars[] = name.toCharArray();
        chars[0] = Character.toLowerCase(chars[0]);
        return new String(chars);
    }
~~~

##### 3.2.4 Java内省机制

**它是 JavaSE 中就有的，对 JavaBean 中属性的默认处理规则**。

回想一下咱写的所有模型类，包括 vo 类，是不是都是写好了属性，之后借助 IDE 生成 `getter` 和 `setter` ，或者借助 `Lombok` 的注解生成 `getter` 和 `setter` ？其实这个生成规则，就是利用了 Java 的内省机制。

**Java 的内省默认规定，所有属性的获取方法以 get 开头（ boolean 类型以 is 开头），属性的设置方法以 set 开头。**根据这个规则，才有的默认的 getter 和 setter 方法。



`Introspector` 类是 Java 内省机制中最核心的类之一，它可以进行很多默认规则的处理（包括获取类属性的 get / set 方法，添加方法描述等），当然它也可以处理这种类名转 beanName 的操作。



## 13. 资源管理

前面介绍到`AbstractXmlApplicationContext` 里面的`loadBeanDefinitions` 方法组合了`XmlBeanDefinitionReader`  进行xml配置解析

~~~java
	protected void loadBeanDefinitions(XmlBeanDefinitionReader reader) throws BeansException, IOException {
		Resource[] configResources = getConfigResources();
		if (configResources != null) {
			reader.loadBeanDefinitions(configResources);
		}
		String[] configLocations = getConfigLocations();
		if (configLocations != null) {
			reader.loadBeanDefinitions(configLocations);
		}
	}
~~~

`Resource` 就是 SpringFramework 中定义的资源模型。

### 1. 原生资源加载

`ClassLoader` 的 `getResource` 和 `getResourceAsStream` 方法，它们本身就是 jdk 内置的加载资源文件的方式。但是spring自己造了一套，原因如下：

jdk 原生的 URL 那一套资源加载方式，对于加载 classpath 或者 `ServletContext` 中的资源来说没有标准的处理手段，而且即便是实现起来也很麻烦。



### 2. spring的资源模型

![image-20230421111223179](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230421111223179.png)

可见最顶级的是`InputStreamSource` 接口

#### 2.1 InputStreamSource

```java
public interface InputStreamSource {
	InputStream getInputStream() throws IOException;
}
```

这个接口只有一个 `getInputStream` 方法，很明显它表达了一件事情：实现了 `InputStreamSource` 接口的实现类，都可以从中取到资源的输入流。

#### 2.2 Resource

然后就是 `InputStreamSource` 的子接口 `Resource` 了

> Interface for a resource descriptor that abstracts from the actual type of underlying resource, such as a file or class path resource.
>
> 它是资源描述符的接口，它可以从基础资源的实际类型中抽象出来，例如**文件或类路径资源**。

这个翻译看起来很生硬，不过咱只需要关注到一个点：**文件或类路径的资源**，仅凭这一个点，咱就可以说，`Resource` 确实更适合 SpringFramework 做资源加载（配置文件通常都放到类路径下）。



#### 2.3 EncodedResource

`EncodedResource` 直接实现了 `InputStreamSource` 接口，从类名上也能看得出来它是编码后的资源。通过源码，发现它内部组合了一个 `Resource` ，说明它本身并不是直接加载资源的。

> EncodedResource 是 Spring Framework 中的一个类，它主要用于表示经过编码的资源，通常是用于加载属性文件或配置文件。**在加载资源文件时，如果文件中包含非 ASCII 字符，则需要对这些字符进行编码以便正确地读取和处理这些字符**。
>
> 通过使用 EncodedResource，Spring 框架可以自动识别并应用正确的编码方式，从而确保在读取资源文件时不会出现乱码等问题。



#### 2.4 WritableResource

`Resource` 有了一个新的子接口：`WritableResource` ，它代表着“可写的资源”，那 `Resource` 就可以理解为“可读的资源”（有木有想起来 `BeanFactory` 与 `ConfigurableBeanFactory` ？）。

#### 2.5 ContextResource

> 从一个封闭的 “上下文” 中加载的资源的扩展接口，例如来自 `javax.servlet.ServletContext` ，也可以来自普通的类路径路径或相对的文件系统路径（在没有显式前缀的情况下指定，因此相对于本地 `ResourceLoader` 的上下文应用）。



### 3. spring的资源模型实现

#### 3.1 Java原生资源加载

大致上分 3 种吧：

- 借助 ClassLoader 加载类路径下的资源

~~~java
        ClassLoader classLoader = MAIN.class.getClassLoader();
        InputStream resourceAsStream = classLoader.getResourceAsStream("Blue.properties");
        Properties properties = new Properties();

        try {
            properties.load(resourceAsStream);
        } catch (Exception e) {
            e.printStackTrace();
        }
        String property = properties.getProperty("blue.order");
        System.out.println(property);
~~~



- 借助 File 加载文件系统中的资源
- 借助 URL 和不同的协议加载本地 / 网络上的资源



#### 3.2 spring的实现

- ClassLoader → `ClassPathResource` [ classpath:/ ]
- File → `FileSystemResource` [ file:/ ]
- URL → `UrlResource` [ xxx:/ ]

除了这三种实现，还有对应于 `ContextResource` 的实现：`ServletContextResource` ，它意味着资源是去 `ServletContext` 域中寻找。



### 4.spring加载资源的方式

提过在 `AbstractApplicationContext` 中，通过类继承关系可以得知它继承了 `DefaultResourceLoader` ，也就是说，**`ApplicationContext` 具有加载资源的能力**。

#### 4.1 DefaultResourceLoader组合了一堆ProtocolResolver

协议解析器

~~~java
private final Set<ProtocolResolver> protocolResolvers = new LinkedHashSet<>(4);

public Resource getResource(String location) {
    Assert.notNull(location, "Location must not be null");

    for (ProtocolResolver protocolResolver : getProtocolResolvers()) {
        Resource resource = protocolResolver.resolve(location, this);
        if (resource != null) {
            return resource;
        }
    }
    // ......
}
~~~

##### 4.1.1 ProtocolResolver

~~~java
// @since 4.3
@FunctionalInterface
public interface ProtocolResolver {
	Resource resolve(String location, ResourceLoader resourceLoader);
}
~~~

它只有一个接口，而且是在 SpringFramework 4.3 版本才出现的（蛮年轻哦），它本身可以搭配 `ResourceLoader` ，在 `ApplicationContext` 中实现**自定义协议的资源加载**，但它还可以脱离 `ApplicationContext` ，直接跟 `ResourceLoader` 搭配即可。这个特性蛮有趣的，咱可以稍微写点代码演示一下效果。



##### 4.1.2 ProtocolResolver使用方式

在工程的 `resources` 目录下新建一个 `Dog.txt` 文件（随便放哪儿都行，只要能找得到），然后写一个 `DogProtocolResolver` ，实现 `ProtocolResolver` 接口：

~~~java
public class DogProtocolResolver implements ProtocolResolver {

    public static final String DOG_PATH_PREFIX = "dog:";

    @Override
    public Resource resolve(String location, ResourceLoader resourceLoader) {
        System.out.println(location);
        if (!location.startsWith(DOG_PATH_PREFIX)) {
            return null;
        }
        // 把自定义前缀去掉
        String realpath = location.substring(DOG_PATH_PREFIX.length());
        String classpathLocation = "classpath:" + realpath;
        return resourceLoader.getResource(classpathLocation);
    }
}

~~~

编写启动类

~~~java
    public static void main(String[] args) throws Exception {
        DefaultResourceLoader defaultResourceLoader = new DefaultResourceLoader();
        //添加一个解析器
        DogProtocolResolver dogProtocolResolver = new DogProtocolResolver();
        defaultResourceLoader.addProtocolResolver(dogProtocolResolver);

        //读取资源
        Resource resource = defaultResourceLoader.getResource("dog:Dog.txt");
        InputStream inputStream = resource.getInputStream();

        //转输入流
        InputStreamReader reader = new InputStreamReader(inputStream, StandardCharsets.UTF_8);
        BufferedReader bufferedReader = new BufferedReader(reader);
        String readLine;
        while ((readLine = bufferedReader.readLine()) != null) {
            System.out.println(readLine);
        }
        bufferedReader.close();
    }
/*
你好
*/
~~~



#### 4.2 DefaultResourceLoader可自行加载类路径下的资源

查看源码

~~~java
public Resource getResource(String location) {
    // ......
    if (location.startsWith("/")) {
        return getResourceByPath(location);
    } else if (location.startsWith(CLASSPATH_URL_PREFIX)) {
        return new ClassPathResource(location.substring(CLASSPATH_URL_PREFIX.length()), getClassLoader());
    }
    // ......
}
~~~

`getResourceByPath`

~~~java
protected Resource getResourceByPath(String path) {
    return new ClassPathContextResource(path, getClassLoader());
}
~~~

不过这个不是绝对的，如果小伙伴现在手头的工程还有引入 `spring-web` 模块的 pom 依赖，会发现 `DefaultResourceLoader` 的几个 Web 级子类中有重写这个方法，以 `GenericWebApplicationContext` 为例：

```java
protected Resource getResourceByPath(String path) {
    Assert.state(this.servletContext != null, "No ServletContext available");
    return new ServletContextResource(this.servletContext, path);
}
```

可以发现这里创建的不再是类路径下了，Web 环境下 SpringFramework 更倾向于从 `ServletContext` 中加载。



#### 4.3  DefaultResourceLoader可支持特定协议

~~~java
public Resource getResource(String location) {
    // ......
    else {
        try {
            // Try to parse the location as a URL...
            URL url = new URL(location);
            return (ResourceUtils.isFileURL(url) ? new FileUrlResource(url) : new UrlResource(url));
        }
        catch (MalformedURLException ex) {
            // No URL -> resolve as resource path.
            return getResourceByPath(location);
        }
    }
}
~~~

如果上面它不能处理类路径的文件，就会尝试通过 URL 的方式加载，这里面包含文件系统的资源，和特殊协议的资源。

所以修改启动类

~~~java
    DogProtocolResolver dogProtocolResolver = new DogProtocolResolver();
//        defaultResourceLoader.addProtocolResolver(dogProtocolResolver);

    //读取资源
    Resource resource = defaultResourceLoader.getResource("Dog.txt");

/*
你好
*/
~~~

也会自动读取resource下的Dog.txt



## 14.  PropertySource的使用

### 1. @PropertySource引入properties文件

#### 1.1  声明properties文件

~~~java
jdbc.url=jdbc:mysql://localhost:3306/test
jdbc.driver-class-name=com.mysql.jdbc.Driver
jdbc.username=root
jdbc.password=123456
~~~

#### 1.2 编写配置模型类

~~~java
@Data
@Component
public class JdbcProperties {
    @Value("${jdbc.url}")
    private String url;

    @Value("${jdbc.driver-class-name}")
    private String driverClassName;

    @Value("${jdbc.username}")
    private String username;

    @Value("${jdbc.password}")
    private String password;

}
~~~

#### 1.3 配置类

~~~java
@Configuration
@ComponentScan("org.example.PropertySource.bean")
@PropertySource("classpath:jdbc.properties")
public class JdbcPropertiesConfiguration {
}
~~~

#### 1.4 测试运行

~~~java
public class MAIN {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(JdbcPropertiesConfiguration.class);
        System.out.println(ctx.getBean(JdbcProperties.class));
    }
}
/*
JdbcProperties(url=jdbc:mysql://localhost:3306/test, driverClassName=com.mysql.jdbc.Driver, username=root, password=123456)
*/
~~~



### 2. 引入xml文件

> 指示要加载的属性文件的资源位置。 支持原生 properties 和基于 XML 的属性文件格式。



#### 2.1 声明xml文件

~~~xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE properties SYSTEM "http://java.sun.com/dtd/properties.dtd">
<properties>
    <entry key="xml.jdbc.url">jdbc:mysql://localhost:3306/test</entry>
    <entry key="xml.jdbc.driver-class-name">com.mysql.jdbc.Driver</entry>
    <entry key="xml.jdbc.username">root</entry>
    <entry key="xml.jdbc.password">123456</entry>
</properties>
~~~

sun 当时给出的 Properties 格式的 xml 标准规范写法，必须按照这个格式来，才能解析为 `Properties` 

#### 2.2 编写配置模型类

注意上面xml文件的key是`xml.xx`

~~~java
@Component
@Data
public class JdbcXmlProperty {

    @Value("${xml.jdbc.url}")
    private String url;

    @Value("${xml.jdbc.driver-class-name}")
    private String driverClassName;

    @Value("${xml.jdbc.username}")
    private String username;

    @Value("${xml.jdbc.password}")
    private String password;

}
~~~



#### 2.3 编写配置类

~~~java
@ComponentScan("org.example.PropertySource.bean")
@Configuration
@PropertySource("classpath:jdbc.xml")
public class JdbcXmlConfiguration {
}

~~~



#### 2.4 测试运行

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(JdbcXmlConfiguration.class);
System.out.println(ctx.getBean(JdbcXmlProperty.class));

/*
JdbcXmlProperty(url=jdbc:mysql://localhost:3306/test, driverClassName=com.mysql.jdbc.Driver, username=root, password=123456)
*/
~~~



#### 2.5 xml格式的限制

为什么要用`<properties>` 来设置属性，先思考spring如何加载.properties文件的，是通过jdk原生`Properties`类

##### 2.5.1 解析Properties入口

~~~java
public @interface PropertySource {
    // ......

	/**
	 * Specify a custom {@link PropertySourceFactory}, if any.
	 * <p>By default, a default factory for standard resource files will be used.
	 * @since 4.3
	 */
	Class<? extends PropertySourceFactory> factory() default PropertySourceFactory.class;
}
~~~

它想表达的是,用 `@PropertySource` 注解引入的资源文件需要用什么策略来解析它。默认情况下它只放了一个 `PropertySourceFactory` 在这里，看一眼 `factory` 属性的泛型也能大概猜得出来，`PropertySourceFactory` 应该是一个接口 / 抽象类，它肯定有默认实现的子类，查找得默认的唯一实现`DefaultPropertySourceFactory`

##### 2.5.2  默认的Properties解析工厂

~~~java
public class DefaultPropertySourceFactory implements PropertySourceFactory {

	@Override
	public PropertySource<?> createPropertySource(@Nullable String name, EncodedResource resource) throws IOException {
		return (name != null ? new ResourcePropertySource(name, resource) : new ResourcePropertySource(resource));
	}

}
~~~

默认实现只是new了一个`ResourcePropertySource`，进入源码得有很多构造方法

~~~java
	public ResourcePropertySource(String name, EncodedResource resource) throws IOException {
		super(name, PropertiesLoaderUtils.loadProperties(resource));
		this.resourceName = getNameForResource(resource.getResource());
	}

	/**
	 * Create a PropertySource based on Properties loaded from the given resource.
	 * The name of the PropertySource will be generated based on the
	 * {@link Resource#getDescription() description} of the given resource.
	 */
	public ResourcePropertySource(EncodedResource resource) throws IOException {
		super(getNameForResource(resource.getResource()), PropertiesLoaderUtils.loadProperties(resource));
		this.resourceName = null;
	}
~~~

注意一个方法`loadProperties`

~~~java
	public static Properties loadProperties(EncodedResource resource) throws IOException {
		Properties props = new Properties();
		fillProperties(props, resource);
		return props;
	}
~~~

可以看到最后还是使用到了`Properties` 进行解析，可是问题来了如果用`Properties`解析xml文件？

##### 2.5.3 【拓展】原生jdkProperties解析xml文件

进行`Properties`源码查看有个`loadFromXml`方法

~~~java
public synchronized void loadFromXML(InputStream in)
    throws IOException, InvalidPropertiesFormatException
{
    Objects.requireNonNull(in);
    PropertiesDefaultHandler handler = new PropertiesDefaultHandler();
    handler.load(this, in);
    in.close();
}
~~~

只是这个 xml 的要求，属实有点高，它是 sun 公司在很早之前就制定的一个 xml 表达 properties 的标准：（以下是 dtd 约束文件内容）

```java
<!--
   Copyright 2006 Sun Microsystems, Inc.  All rights reserved.
  -->

<!-- DTD for properties -->

<!ELEMENT properties ( comment?, entry* ) >

<!ATTLIST properties version CDATA #FIXED "1.0">

<!ELEMENT comment (#PCDATA) >

<!ELEMENT entry (#PCDATA) >

<!ATTLIST entry key CDATA #REQUIRED>
```

可以发现确实是有固定格式的，必须按照这个约束来编写 xml 文件。这个东西知道就好了，估计你以后也用不到 ~ ~ ~



##### 2.5.4 properties与xml的对比

```properties
jdbc.url=jdbc:mysql://localhost:3306/test
jdbc.driver-class-name=com.mysql.jdbc.Driver
jdbc.username=root
jdbc.password=123456
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE properties SYSTEM "http://java.sun.com/dtd/properties.dtd">
<properties>
    <entry key="xml.jdbc.url">jdbc:mysql://localhost:3306/test</entry>
    <entry key="xml.jdbc.driver-class-name">com.mysql.jdbc.Driver</entry>
    <entry key="xml.jdbc.username">root</entry>
    <entry key="xml.jdbc.password">123456</entry>
</properties>
```

难易程度高下立判，properties 完胜，所以对于这种配置型的资源文件，通常都是使用 properties 来编写。

当然，properties 也不是完全 OK ，由于它的特征是 key-value 的形式，整个文件排下来是**没有任何层次性可言的**（换句话说，每个配置项之间的地位都是平等的）。这个时候 xml 的优势就体现出来了，它可以非常容易的体现出层次性，不过咱不能因为这一个点就觉得 xml 还可以，因为有一个更适合解决这个问题的配置格式：**yml** 。

### 3. @PropertySource引入yml文件

**yml** 又称 **yaml** ，它是可以代替 properties 同时又可以表达层级关系的标记语言，它的基本格式如下：

~~~yaml
person:
  name: zhangsan
  age: 18
  cat:
    name: mimi
    color: white
dog:
  name: wangwang
~~~

这种写法同等于下面的 properties ：

```properties
person.name=zhangsan
person.age=18
person.cat.name=mimi
person.cat.color=white
dog.name=wangwang
```



#### 3.1 声明yml文件

~~~yaml
yml: 
  jdbc:
    url: jdbc:mysql://localhost:3306/test
    driver-class-name: com.mysql.jdbc.Driver
    username: root
    password: 123456
~~~



#### 3.2 编写模型类

同上面两种，只不过前缀改了

~~~java
@Component
@Data
public class JdbcYmlProperty {

    @Value("${yml.jdbc.url}")
    private String url;

    @Value("${yml.jdbc.driver-class-name}")
    private String driverClassName;

    @Value("${yml.jdbc.username}")
    private String username;

    @Value("${yml.jdbc.password}")
    private String password;

}
~~~



#### 3.3 配置类

~~~java
@ComponentScan("org.example.PropertySource.bean")
@Configuration
@PropertySource("classpath:jdbc.yml")
public class JdbcYmlConfiguration {
}
~~~



#### 3.4 测试运行

~~~sh
JdbcYmlProperty(url=${yml.jdbc.url}, driverClassName=${yml.jdbc.driver-class-name}, username=${yml.jdbc.username}, password=${yml.jdbc.password})
~~~

发现配置属性并没有注入，因为`@PropertySource` 默认的解析策略是`*DefaultPropertySourceFactory*`,只会解析`xml`和`properties` 文件，所以无法解析`yml`文件，只能自己设计。

#### 3.5 自定义yml解析

##### 3.5.1 导入依赖

使用snakeyaml库进行解析

~~~xml
<dependency>
    <groupId>org.yaml</groupId>
    <artifactId>snakeyaml</artifactId>
    <version>1.26</version>
</dependency>
~~~

##### 3.5.2 自定义PropertySourceFactory

~~~java
public class YamlPropertySourceFactory implements PropertySourceFactory {
    @Override
    public PropertySource<?> createPropertySource(String name, EncodedResource resource) throws IOException {        
        return null;
    }
}
~~~

注意看这个接口的方法，它要返回一个 `PropertySource<?>` ，借助 IDEA 观察它的继承关系，可以发现它里头有一个实现类叫 `PropertiesPropertySource` 

![image-20230421154115100](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230421154115100.png)

查看源码得，只有一个公开构造方法，所以我们就要将yml资源文件对象转为Properties对象即可

~~~java
public PropertiesPropertySource(String name, Properties source) {
    super(name, (Map) source);
}

protected PropertiesPropertySource(String name, Map<String, Object> source) {
    super(name, source);
}
~~~



##### 3.5.3 资源文件转Properties对象

~~~java
public class YamlPropertySourceFactory implements PropertySourceFactory {
    @Override
    public PropertySource<?> createPropertySource(String name, EncodedResource resource) throws IOException {
        // 获得yaml配置
        YamlPropertiesFactoryBean yamlPropertiesFactoryBean = new YamlPropertiesFactoryBean();
        yamlPropertiesFactoryBean.setResources(resource.getResource());
        
        //直接解析获得Properties对象
        Properties properties = yamlPropertiesFactoryBean.getObject();

        assert properties != null;
        return new PropertiesPropertySource(name != null ? name : Objects.requireNonNull(resource.getResource().getFilename()), properties);
    }
}
~~~



##### 3.5.4 设置`@PropertySource`

把这个 `YmlPropertySourceFactory` 设置到 `@PropertySource` 中：

~~~java
@PropertySource(value = "classpath:jdbc.yml",factory = YamlPropertySourceFactory.class)
~~~



##### 3.5.5 测试运行得

~~~sh
JdbcYmlProperty(url=jdbc:mysql://localhost:3306/test, driverClassName=com.mysql.jdbc.Driver, username=root, password=123456)
~~~



**@PropertySource默认只能加载properties和特定格式的xml文件，最终默认实现都是返回JDK自带的Properties类，该注解提供自定义工厂来解析处理**



## 15. 配置源&配置元信息

### 1.配置源

**咱现在已经学过的配置源就是 xml 配置文件，以及注解配置类两种**

#### 1.1 如何理解配置源

**配置源**，简单理解，就是**配置的来源**。在前面的超多例子中，都是使用 xml 配置文件或者注解配置类来驱动 IOC 容器，那么对于 IOC 容器而言，xml 配置文件或者注解配置类就可以称之为配置源。

#### 1.2 配置源的解析思路

##### 1.2.1 xml配置文件

~~~xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:context="http://www.springframework.org/schema/context"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd http://www.springframework.org/schema/context https://www.springframework.org/schema/context/spring-context.xsd">

    <context:component-scan base-package="com.linkedbear.spring.basic_di.c_value_spel.bean"/>

    <context:property-placeholder location="classpath:basic_di/value/red.properties"/>

    <bean id="person" class="com.linkedbear.spring.basic_di.a_quickstart_set.bean.Person">
        <property name="name" value="test-person-byset"/>
        <property name="age" value="18"/>
    </bean>
</beans>
~~~

这个 xml 中包含几个部分：

- xml 头信息
- `component-scan` 声明包扫描
- `property-placeholder` 引入外部 properties 文件
- `<bean>` 注册 bean 并属性赋值

这个 xml 可以用如下的一种抽象语言描述：

```ini
beans.xml {
    context: [component-scan, property-placeholder]
    beans: [person]
}
```

这里面不会描述具体的组件扫描路径等等，只会**记录这个 xml 中声明了哪些标签**。

##### 1.2.2 注解配置类

跟上面一样，咱先找一个之前写过的注解配置类：

```java
@Configuration
@ComponentScan("com.linkedbear.spring.bean.b_scope.bean")
public class BeanScopeConfiguration {
    
    @Bean
    public Child child1(Toy toy) {
        Child child = new Child();
        child.setToy(toy);
        return child;
    }
    
    @Bean
    public Child child2(Toy toy) {
        Child child = new Child();
        child.setToy(toy);
        return child;
    }
}
```

根据上面的抽象思维，这个注解配置类也可以进行如下转换：

```yaml
BeanScopeConfiguration.java: {
    annotations: [ComponentScan，Configuration]
    beans: [child1, child2]
}
```

与上面一样，只会记录配置类中的**配置结构**而已，任何配置信息都不会体现在这里面。

### 2. 元信息

#### 2.1 如何理解元信息

**元信息**，又可以理解为**元定义**，简单的说，它就是**定义的定义**。

干说太抽象，直接上例子吧，方便理解：

- 张三，男，18岁
  - 它的元信息就是它的属性们：`Person {name, age, sex}`
- 咪咪，美国短毛，黑白毛，主人是张三
  - 它的元信息可以抽取为：`Cat {name, type, color, master}`

写到这里是不是突然产生一种感觉：这不就是**对象和类**吗？？？是的，所以我们可以这样说：**类中包含对象的元信息**。

类有元信息吗？当然有，**`Class`** 这个类里面就包含一个类的所有定义（属性、方法、继承实现、注解等），所以我们可以说：**`Class` 中包含类的元信息**。

**数据库表与表结构信息，这也是非常典型的信息与元信息：数据库表结构描述了数据库表的整体表属性，以及表字段的属性。**

#### 2.2 spring中配置元信息

##### 2.2.1 Bean的定义元信息

SpringFramework 中定义的 Bean 也会封装为一个个的 Bean 的元信息，也就是 **`BeanDefinition`**

- Bean 的全限定名 className
- Bean 的作用域 scope
- Bean 是否延迟加载 lazy
- Bean 的工厂 Bean 名称 factoryBean
- Bean 的构造方法参数列表 constructorArgumentValues
- Bean 的属性值 propertyValues
- ......



##### 2.2.2 IOC容器的配置元信息

**IOC 容器的配置元信息分为 beans 和 context 两部分**



###### 2.2.2.1 beans的配置元信息

以 xml 配置文件为例，如果你仔细注意一下整个配置文件的最顶层标签，会发现 `<beans>` 其实是有属性的：

![image-20230423104910542](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230423104910542.png)

但一般不使用，配置了默认值

| 配置元信息                  | 含义 / 作用                                                  | 默认值           |
| --------------------------- | ------------------------------------------------------------ | ---------------- |
| profile                     | 基于环境的配置                                               | ""               |
| default-autowire            | 默认的自动注入模式（不需要声明 `@Autowired` 等注解即可注入组件） | default（no）    |
| default-autowire-candidates | 满足指定属性名规则的属性才会被自动注入                       | ""               |
| default-init-method         | 全局 bean 的初始化方法                                       | ""               |
| default-destroy-method      | 全局 bean 的销毁方法                                         | ""               |
| default-lazy-init           | 全局 bean 是否延迟加载                                       | default（false） |
| default-merge               | 继承父 bean 时直接合并父 bean 的属性值                       | default（false） |

> 默认值中提到的 default 是在没有声明时继承父配置的默认值（ `<beans>` 标签是可以嵌套使用的），如果都没有声明，则配置的默认值是括号内的值。



**其他配置元信息**

beans的命名空间还有两个常见标签



| 配置元信息   | 含义 / 作用                 | 使用方式举例                                                 |
| ------------ | --------------------------- | ------------------------------------------------------------ |
| `<alias />`  | 给指定的 bean 指定别名      | `<alias name="person" alias="zhangsan"/>`                    |
| `<import />` | 导入外部现有的 xml 配置文件 | `<import resource="classpath:basic_dl/quickstart-byname.xml"/>` |



###### 2.2.2.2 context的配置元信息

![image-20230423105231590](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230423105231590.png)

![image-20230423105332859](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230423105332859.png)

常用的是前三个配置。



##### 2.2.3 properties等配置元信息

我们反复研究的 properties 、xml 、yml 文件，它们的作用都是为了将具体的配置抽取为一个可任意修改的配置文件，防止在程序代码中出现硬编码配置，导致修改配置还需要重新编译的麻烦。这种**将配置内容抽取为配置文件**的动作，我们称之为 **“配置外部化”**，抽取出来的配置文件又被成为 **“外部化配置文件”** 

> 而加载这些外部化配置文件的方式，要么通过上面的 `<context:property-placeholder/>` ，要么通过 `@PropertySource` 注解，它们最终都会被封装为一个一个的 `PropertySource` 对象（ properties 文件被封装为 `PropertiesPropertySource` ）了，而这个 `PropertySource` 对象内部就持有了这些外部化配置文件的所有内容。



## 16. Environment抽象

> 加载的 properties 资源配置，以及 `ApplicationContext` 内部的一些默认配置属性，都放在哪里了？组件 Bean 又是怎么把配置值注入进去到对象的属性中的

### 1. Environment概述

#### 1.1 第一感受

其实第一眼看到这个名词，我们就应该有一个模糊的猜想了，它应该是基于 SpringFramework 的工程的**运行时环境**。

![img](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/461c0349a0964aa8925dc92648796803~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

#### 1.2 官方文档

> `Environment` 接口是集成在容器中的抽象，可对应用程序环境的两个关键方面进行建模：Profile 和 properties 。 Profiles 是仅在指定 profile 处于活动状态（ active ）时才向容器注册 `BeanDefinition` 的命名逻辑组。它可以将 Bean 分配给不同的 profile （无论是以 XML 定义还是注解配置）。与配置文件相关的 `Environment` 作用是确定哪些配置文件当前处于活动状态，以及哪些配置文件在默认情况下应处于活动状态。 Properties 在几乎所有应用程序中都起着重要作用，并且可能源自多种来源：属性文件，JVM 系统属性，系统环境变量，JNDI，`ServletContext` 参数，临时属性对象，`Map` 对象等。`Environment` 与属性相关联的作用是为用户提供方便的接口，它可以用于配置属性源，并从 `Environment` 中解析属性。

意思是运行时环境可以限制配置文件或者属性运行的时刻。

第一句【`Environment` 是集成在容器中的抽象】，会让我们产生一种感觉：前面的理解是不是出现了一些偏差，那么实际模型就该这样：

![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/5166a2a211cf4cb1a36a6bb4f5232755~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



#### 1.3 见解

- 首先，`Environment` 中包含 profiles 和 properties ，这些配置信息会影响 IOC 容器中的 bean 的注册与创建；
- 其次，`Environment` 的创建是在 `ApplicationContext` 创建后才创建的（ IOC 原理部分会解释），所以 `Environment` 应该是伴随着 `ApplicationContext` 的存在而存在；
- 第三，`ApplicationContext` 中同时包含 `Environment` 和组件 bean ，而且从 `BeanFactory` 的视角来看，`Environment` 也是一个 Bean ，只不过它的地位比较特殊。

所以工程结构应该如下

![img](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/07df885666e94a248600430502779905~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

#### 1.4 描述

##### 1.4.1 Environment包含profile与properties

> `Environment` 是表示当前应用程序正在其中运行的环境的接口。它为应用环境制定了两个关键的方面：**profile** 和 **properties**。与属性访问有关的方法通过 `PropertyResolver` 这个父接口公开。

`PropertyResolver` 这个接口，这个接口负责解析占位符（ **${...}** ）对应的值

##### 1.4.2 profile用于区分不同的环境模式

>  profile 机制保证了仅在给定 profile 处于激活状态时，才向容器注册的 `BeanDefinition` 的命名逻辑组.

`Environment` 配合 profile 可以完成**指定模式的环境的组件装配**，以及不同的配置属性注入。

##### 1.4.3 properties用于配置属性和注入值

 properties 的最大作用之一是做**外部化配置**，`Environment` 中存放了很多 properties ，它们的来源有很多种，而最终的作用都是**提供了属性配置**，或者**给组件注入属性值**。

##### 1.4.4 Environment不建议直接使用

> 在 `ApplicationContext` 中管理的 Bean 可以注册为 `EnvironmentAware` 或使用 `@Inject` 标注在 `Environment` 上，以便直接查询 profile 的状态或解析 `Properties`。 但是，在大多数情况下，应用程序级 Bean 不必直接与 `Environment` 交互，而是通过将 **${...}** 属性值替换为属性占位符配置器进行属性注入（例如 `PropertySourcesPlaceholderConfigurer`），该属性本身是 `EnvironmentAware`，当配置了 `<context:property-placeholder/>` 时，默认情况下会使用 Spring 3.1 的规范注册。

`Environment` 可以注入到组件中，用于获取当前环境激活的所有 profile 模式；但是又不推荐开发者直接使用它，而是通过占位符注入配置属性的值。为什么会这么说呢，其实这个又要说回 `Environment` 设计的原始意图。`Environment` 的设计本身就应该是一个**不被应用程序接触到的 “环境”** ，我们**只能从环境中获取一些它已经有的信息，但不应该获取它本身**。所以，在处理 properties 的获取时，直接使用占位符就可以获取了。

`Blue.properties`

~~~java
user.name=123
user.id=12
~~~

测试`PropertySourcesPlaceholderConfigurer`和`PropertyPlaceholderConfigurer`

~~~java
public class EnvironmentPlaceholderDemo {

    @Value("${user.id}")
    private Long id;

    @Value("${user.name}")
    private String name;

    public static void main(String[] args) {
        AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext();
        // 注册 Configuration Class
        context.register(EnvironmentPlaceholderDemo.class);

        // 启动 Spring 应用上下文
        context.refresh();

        EnvironmentPlaceholderDemo environmentPlaceholderDemo = context.getBean(EnvironmentPlaceholderDemo.class);

        System.out.println(environmentPlaceholderDemo.id);
        System.out.println(environmentPlaceholderDemo.name);

        // 关闭 Spring 应用上下文
        context.close();
    }

    /**
     * Spring 3.1前使用PropertyPlaceholderConfigurer处理占位符
     * 加上static保证提前初始化
     * user.name = 123
     */
//    @Bean
//    public static PropertyPlaceholderConfigurer propertyPlaceholderConfigurer() {
//        PropertyPlaceholderConfigurer propertyPlaceholderConfigurer = new PropertyPlaceholderConfigurer();
//        propertyPlaceholderConfigurer.setLocation(new ClassPathResource("Blue.properties"));
//        propertyPlaceholderConfigurer.setFileEncoding("UTF-8");
//        return propertyPlaceholderConfigurer;
//    }

    /**
     * Spring 3.1 + 推荐使用PropertySourcesPlaceholderConfigurer
     * 这里的user.name显示的并不是张三，而是电脑的name，这涉及外部化配置
     */
//    @Bean
//    public static PropertySourcesPlaceholderConfigurer propertySourcesPlaceholderConfigurer() {
//        PropertySourcesPlaceholderConfigurer propertySourcesPlaceholderConfigurer = new PropertySourcesPlaceholderConfigurer();
//        propertySourcesPlaceholderConfigurer.setLocation(new ClassPathResource("Blue.properties"));
//        propertySourcesPlaceholderConfigurer.setFileEncoding("UTF-8");
//        return propertySourcesPlaceholderConfigurer;
//    }

}
~~~



##### 1.4.5 ApplicationContext获取到的是ConfigurableEnvironment

> `ApplicationContext` 的根实现类 `AbstractApplicationContext` 获取到的是 `ConfigurableEnvironment` ，它具有 **“可写”** 的特征，换言之我们可以修改它内部的属性值 / 数据。不过话又说回来，通常情况下我们都不会直接改它，除非要对 SpringFramework 应用的启动流程或者运行中进行一些额外的扩展或者修改。

~~~java
/** Environment used by this context. */
@Nullable
private ConfigurableEnvironment environment;
~~~

#### 1.5 面试介绍

**`Environment` 是 SpringFramework 3.1 引入的抽象的概念，它包含 profiles 和 properties 的信息，可以实现统一的配置存储和注入、配置属性的解析等。其中 profiles 实现了一种基于模式的环境配置，properties 则应用于外部化配置。**

### 2. Environment的结构

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/7ace9b0f0c144ca4a49e3083fbece4eb~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

#### 2.1 PropertyResolver

这个接口，只从接口名就知道它应该是处理占位符 **${}** 的。观察接口的方法定义，直接实锤了它就是做配置属性值的获取和解析的：（下面是 `PropertyResolver` 的部分方法定义）

```java
public interface PropertyResolver {

    // 检查所有的配置属性中是否包含指定key
    boolean containsProperty(String key);

    // 以String的形式返回指定的配置属性的值
    String getProperty(String key);

    // 带默认值的获取
    String getProperty(String key, String defaultValue);

    // 指定返回类型的配置属性值获取
    <T> T getProperty(String key, Class<T> targetType);

    // ......

    // 解析占位符
    String resolvePlaceholders(String text);

    // ......
}
```

所以由此也就证明了：**`Environment` 可以获取配置元信息，同时也可以解析占位符的信息**。



#### 2.2 ConfigurableEnvironment

看到 **Configurable** 开头，就知道它是**可配置的**类。

~~~java
void setActiveProfiles(String... profiles);

void addActiveProfile(String profile);

void setDefaultProfiles(String... profiles);

MutablePropertySources getPropertySources();
~~~

注意到一个`MutablePropertySources` 

~~~java
public class MutablePropertySources implements PropertySources {

	private final List<PropertySource<?>> propertySourceList = new CopyOnWriteArrayList<>();
~~~

获得资源集合。可得：**Mutable 开头的类名，通常可能是一个类型的 List 组合封装**。

#### 2.3 *StandardEnvironment*

~~~java
	/** System environment property source name: {@value}. */
	public static final String SYSTEM_ENVIRONMENT_PROPERTY_SOURCE_NAME = "systemEnvironment";

	/** JVM system properties property source name: {@value}. */
	public static final String SYSTEM_PROPERTIES_PROPERTY_SOURCE_NAME = "systemProperties";
~~~

在`AbstractEnvironment` 基础上新增两个属性源，按照`system properties
system environment variables` 进行搜寻。



### 3. Environment的基本使用

虽说 `Environment` 不建议直接在应用程序中使用，但是部分场景下还是需要直接接触它来操纵。

#### 3.1 获得Environment的API

既然 `Environment` 存在于 `ApplicationContext` 中，那么获取 `Environment` 的方式自然也就可以想到：`@Autowired`

~~~java
@Component
public class EnvironmentHolder {

    @Autowired
    private Environment environment;

    public void PrintEnvironment() {
        System.out.println(environment);
    }

}
~~~

测试得

~~~java
public class MAIN {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext("org.example.PropertySource.bean");
        ctx.getBean(EnvironmentHolder.class).PrintEnvironment();
    }
}
/*
StandardEnvironment {activeProfiles=[], defaultProfiles=[default], propertySources=[PropertiesPropertySource {name='systemProperties'}, SystemEnvironmentPropertySource {name='systemEnvironment'}]}
*/
~~~



除此之外，联想到 `BeanFactory` 、`ApplicationContext` 的注入方式还有回调注入，作为 SpringFramework 的内置 API ，估计也会有一个 **Aware** 回调注入的接口吧！那自然是必须的，`EnvironmentAware` 就是回调注入的接口

~~~java
@Component
public class EnvironmentHolder implements EnvironmentAware {

    private Environment environment;

    public void PrintEnvironment() {
        System.out.println(environment);
    }

    @Override
    public void setEnvironment(Environment environment) {
        this.environment = environment;
    }
}
~~~

> 注：使用 `@Autowired` 的方式在某些情况下会注入失败，所以对于小伙伴们而言，注入是否能成功需要亲手测试运行检验才能知道。在后面的后置处理器部分，会演示一种无法使用 `@Autowired` 注入 `Environment` 的方式，小伙伴们到时候可以留意一下。



#### 3.2 使用Environment获取配置属性的值

配置类

~~~java
@Configuration
@ComponentScan("org.example.PropertySource.bean")
@PropertySource("Blue.properties")
public class EnvironmentPropertyConfiguration {
}
~~~

bean

~~~java
public void PrintEnvironment() {
    System.out.println(environment.getProperty("user.name"));
    System.out.println(Arrays.toString(environment.getDefaultProfiles()));
}
~~~

启动测试得

~~~sh
wuxie #本机的username
[default]
~~~



### 4.原理探讨

#### 4.1  默认Profiles

可以发现默认的环境是`default`，从何而来就进入他的抽象实现类`AbstractEnvironment`

~~~java
protected static final String RESERVED_DEFAULT_PROFILE_NAME = "default";
@Override
public String[] getDefaultProfiles() {
    return StringUtils.toStringArray(doGetDefaultProfiles());
}
~~~

`doGetDefaultProfiles`，为该方法前加了个`do`，这个设计很常见

##### 4.1 .1 方法命名规范

在 SpringFramework 的框架编码中，如果有出现一个方法是 do 开头，并且去掉 do 后能找到一个与剩余名称一样的方法，则代表如下含义：**不带 do 开头的方法一般负责前置校验处理、返回结果封装，带 do 开头的方法是真正执行逻辑的方法（如 `getBean` 方法的底层会调用 `doGetBean` 来真正的寻找 IOC 容器的 bean ，`createBean` 会调用 `doCreateBean` 来真正的创建一个 bean ）。**



##### 4.1.2 doGetDefaultProfiles的实现

~~~java
protected Set<String> doGetDefaultProfiles() {
    synchronized (this.defaultProfiles) {
        // 取默认的profiles及逆行对比
        if (this.defaultProfiles.equals(getReservedDefaultProfiles())) {
            // 如果一致就获得profiles
            String profiles = doGetDefaultProfilesProperty();
            
            //如果有，就覆盖掉原来默认值
            if (StringUtils.hasText(profiles)) {
                setDefaultProfiles(StringUtils.commaDelimitedListToStringArray(
                        StringUtils.trimAllWhitespace(profiles)));
            }
        }
        return this.defaultProfiles;
    }
}
~~~



~~~java
public static final String DEFAULT_PROFILES_PROPERTY_NAME = "spring.profiles.default";
=======    
@Nullable
protected String doGetDefaultProfilesProperty() {
    return getProperty(DEFAULT_PROFILES_PROPERTY_NAME);
}
~~~



##### 4.1.3 覆盖默认的profiles方法

声明式进行设置defaultProfiles

![image-20230423144628984](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230423144628984.png)



#### 4.2 Environment解析properties

已经知道 `Environment` 继承了父接口 `PropertyResolver` ，自然它拥有解析配置元信息的能力，但如何实现呢

##### 4.2.1  PropertyResolver的实现类

![image-20230423145018150](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230423145018150.png)

查看最终的实现类`StandardEnvironment`

##### 4.2.2 getProperty实现是委派

进入了`StandardEnvironment` 发现没有获得属性的方法。都在父类`AbstractEnvironment`里，并且也不算自己实现，而是调用了一个解析器。

~~~java
public String getProperty(String key) {
    return this.propertyResolver.getProperty(key);
}
...
    
	private final ConfigurablePropertyResolver propertyResolver;    
~~~

一般的，我们称这种方式叫做 **“委派”** ，它与代理、装饰者不同：**委派仅仅是将方法的执行转移给另一个对象，而代理可能会在此做额外的处理，装饰者也会在方法执行前后做增强**。

ConfigurablePropertyResolver最后实现类也是`PropertySourcesPropertyResolver`。

所以总的来说：`Environment` 的解析配置属性值的底层是交给 `PropertySourcesPropertyResolver` 来处理



## 17.Bean与BeanDefinition

### 1.BeanDefinition

#### 1.1 官方介绍

`BeanDefinition` 也是一种**配置元信息**，它描述了 **Bean 的定义信息**。

> bean 的定义信息可以包含许多配置信息，包括构造函数参数，属性值和特定于容器的信息，例如初始化方法，静态工厂方法名称等。**子 bean 定义可以从父 bean 定义继承配置数据**。子 bean 的定义信息可以覆盖某些值，或者可以根据需要添加其他值。使用父 bean 和子 bean 的定义可以节省很多输入（实际上，这是一种模板的设计形式）。



#### 1.2 javadoc

> `BeanDefinition` 描述了一个 bean 的实例，该实例具有属性值，构造函数参数值以及具体实现所提供的更多信息。 这只是一个最小的接口，它的主要目的是允许 `BeanFactoryPostProcessor`（例如 `PropertyPlaceholderConfigurer` ）内省和修改属性值和其他 bean 的元数据。

javadoc 额外提了编码设计中 `BeanDefinition` 的使用：`BeanFactoryPostProcessor` 可以任意修改 `BeanDefinition` 中的信息



#### 1.3 BeanDefinition接口的方法定义

`BeanDefinition` 整体包含以下几个部分：

- Bean 的类信息 - 全限定类名 ( beanClassName )
- Bean 的属性 - 作用域 ( scope ) 、是否默认 Bean ( primary ) 、描述信息 ( description ) 等
- Bean 的行为特征 - 是否延迟加载 ( lazy ) 、是否自动注入 ( autowireCandidate ) 、初始化 / 销毁方法 ( initMethod / destroyMethod ) 等
- Bean 与其他 Bean 的关系 - 父 Bean 名 ( parentName ) 、依赖的 Bean ( dependsOn ) 等
- Bean 的配置属性 - 构造器参数 ( constructorArgumentValues ) 、属性变量值 ( propertyValues ) 等

由此可见，`BeanDefinition` 几乎把 bean 的所有信息都能收集并封装起来，可以说是很全面了。



#### 1.4 [面试题] 如何概述BeanDefinition

**`BeanDefinition` 描述了 SpringFramework 中 bean 的元信息，它包含 bean 的类信息、属性、行为、依赖关系、配置信息等。`BeanDefinition` 具有层次性，并且可以在 IOC 容器初始化阶段被 `BeanDefinitionRegistryPostProcessor` 构造和注册，被 `BeanFactoryPostProcessor` 拦截修改等。**



### 2.BeanDefinition结构

![img](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/3e4f5db098034b6385f9645092990679~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)



#### 2.1 AttributeAccessor

**属性的访问器**

> 定义用于将元数据**附加到任意对象**，或从任意对象**访问元数据**的通用协定的接口。

它很像是类中定义的 getter 和 setter 方法呀，不过实际上它与 getter 、setter 方法有所区别。

##### 2.1.1 回顾元信息概念

~~~java
public class Person {
    private String name;
    
    public String getName() {
        return name;
    }
    
    public void setName(String name) {
        this.name = name;
    }
}
~~~

抽象如下内容

~~~java
className: Person
attributes: [name]
methods: [getName, setName]
~~~



##### 2.1.2 AttributeAccessor 设计

~~~java
public interface AttributeAccessor {
    // 设置bean中属性的值
    void setAttribute(String name, @Nullable Object value);

    // 获取bean中指定属性的值
    Object getAttribute(String name);

    // 移除bean中的属性
    Object removeAttribute(String name);

    // 判断bean中是否存在指定的属性
    boolean hasAttribute(String name);

    // 获取bean的所有属性
    String[] attributeNames();
}
~~~

总结出第一个 `BeanDefinition` 的特征：**`BeanDefinition` 继承了 `AttributeAccessor` 接口，具有配置 bean 属性的功能。**

#### 2.2 BeanMetadataElement

接口名中`Metadata` ，可见是存放了**bean的元信息**，接口只有一个

~~~java
@Nullable
default Object getSource() {
    return null;
}
~~~

> 返回此元数据元素的配置源Object，就是 bean 的文件 / url 路径



#### 2.3 AbstractBeanDefinition

`BeanDefinition` 的第一个实现类，里面已经定义好了一些属性和功能。



 ~~~java
    // bean的全限定类名
     private volatile Object beanClass;
 
     // 默认的作用域为单实例
     private String scope = SCOPE_DEFAULT;
 
     // 默认bean都不是抽象的
     private boolean abstractFlag = false;
 
     // 是否延迟初始化
     private Boolean lazyInit;
     
     // 自动注入模式(默认不自动注入)
     private int autowireMode = AUTOWIRE_NO;
 
     // 是否参与IOC容器的自动注入(设置为false则它不会注入到其他bean，但其他bean可以注入到它本身)
     // 可以这样理解：设置为false后，你们不要来找我，但我可以去找你们
     private boolean autowireCandidate = true;
 
     // 同类型的首选bean
     private boolean primary = false;
 
     // bean的构造器参数和参数值列表
     private ConstructorArgumentValues constructorArgumentValues;
 
     // bean的属性和属性值集合
     private MutablePropertyValues propertyValues;
 
     // bean的初始化方法
     private String initMethodName;
 
     // bean的销毁方法
     private String destroyMethodName;
 
     // bean的资源来源
     private Resource resource;
 ~~~

可见方法几乎很全，但是为什么还要抽象出来？根据文档注册得

> 它是 `BeanDefinition` 接口的抽象实现类，其中排除了 `GenericBeanDefinition` ，`RootBeanDefinition` 和 `ChildBeanDefinition` 的常用属性。 自动装配常量与 `AutowireCapableBeanFactory` 接口中定义的常量匹配。

所以针对不同的 `BeanDefinition` 落地实现，还有一些特殊的属性咯，所以还是需要抽象出一个父类才行哈。

里面有个属性`autowireMode`

##### 2.3.1 补充：自动注入模式

前面讲到bean的配置元信息得xml文件里面有个属性`default-autowire(默认no)`

正常来讲，bean 中的组件依赖注入，是需要在 xml 配置文件，或者在属性 / 构造器 / setter 方法上标注注入的注解（ `@Autowired` / `@Resource` / `@Inject` 的。不过，SpringFramework 为我们提供了另外一种方式，**如果组件中的类型 / 属性名与需要注入的 bean 的类型 / name 完全一致，可以不标注依赖注入的注解，也能实现依赖注入**。



一般情况下，自动注入只会在 xml 配置文件中出现，注解配置中 `@Bean` 注解的 `autowire` 属性在 SpringFramework 5.1 之后被标记为已过时，替代方案是使用 `@Autowired` 等注解。

使用方式则是在xml的bean标签上声明注册模式，`byName` 根据名称注入即可

~~~xml
<bean id="cat" class="com.linkedbear.spring.basic_di.a_quickstart_set.bean.Cat" autowire="byName">
    <property name="name" value="test-cat"/>
    <!-- <property name="master" ref="person"/> 可以不写 -->
</bean>
~~~

自动注入的模式有 5 种选择：`AUTOWIRE_NO`（不自动注入）、`AUTOWIRE_BY_NAME`（根据 bean 的名称注入）、`AUTOWIRE_BY_TYPE`（根据 bean 的类型注入）、`AUTOWIRE_CONSTRUCTOR`（根据 bean 的构造器注入）、`AUTOWIRE_AUTODETECT`（借助内省决定如何注入，3.0 即弃用），**默认是不开启的**（所以才需要我们开发者对需要注入的属性标注注解，或者在 xml 配置文件中配置）。

#### 2.4 实现类 GenericBeanDefinition

又看到 `Generic` 了，它代表着通用、一般的，所以这种 `BeanDefinition` 也具有一般性。`GenericBeanDefinition` 的源码实现非常简单，仅仅是比 `AbstractBeanDefinition` 多了一个 `parentName` 属性而已。

由这个设计，可以得出以下几个结论：

- `AbstractBeanDefinition` 已经完全可以构成 `BeanDefinition` 的实现了
- `GenericBeanDefinition` 就是 `AbstractBeanDefinition` 的非抽象扩展而已
- `GenericBeanDefinition` 具有层次性（可从父 `BeanDefinition` 处继承一些属性信息）



#### 2.5 RootBeanDefinition与ChildBeanDefinition

前缀`Root/Child`，字面意思为根/子。

`ChildBeanDefinition` ，它的设计实现与 `GenericBeanDefinition` 如出一辙，都是集成一个 `parentName` 来作为父 `BeanDefinition` 的 “指向引用” 

`RootBeanDefinition` 有着 “根” 的概念在里面，它只能作为单体独立的 `BeanDefinition` ，或者父 `BeanDefinition` 出现



下面是 `RootBeanDefinition` 的一些重要的成员属性：

```java
    // BeanDefinition的引用持有，存放了Bean的别名
    private BeanDefinitionHolder decoratedDefinition;

    // Bean上面的注解信息
    private AnnotatedElement qualifiedElement;

    // Bean中的泛型
    volatile ResolvableType targetType;

    // BeanDefinition对应的真实的Bean
    volatile Class<?> resolvedTargetType;

    // 是否是FactoryBean
    volatile Boolean isFactoryBean;
    // 工厂Bean方法返回的类型
    volatile ResolvableType factoryMethodReturnType;
    // 工厂Bean对应的方法引用
    volatile Method factoryMethodToIntrospect;
```

可以发现，`RootBeanDefinition` 在 `AbstractBeanDefinition` 的基础上，又扩展了这么些 Bean 的信息：

- Bean 的 id 和别名
- Bean 的注解信息
- Bean 的工厂相关信息（是否为工厂 Bean 、工厂类、工厂方法等）

#### 2.6 AnnotatedBeanDefinition

~~~java
public interface AnnotatedBeanDefinition extends BeanDefinition {
    AnnotationMetadata getMetadata();

    @Nullable
    MethodMetadata getFactoryMethodMetadata();
}
~~~

由这个接口定义的方法，大概就可以猜测到，它可以把 Bean 上的注解信息提供出来。借助 IDEA ，发现它的子类里，有一个 `AnnotatedGenericBeanDefinition` ，还有一个 `ScannedGenericBeanDefinition` ，它们都是基于注解驱动下的 Bean 的注册，封装的 `BeanDefinition` 。

### 3. BeanDefinition使用

#### 3.1 基于xml

使用 xml 配置文件的方式，每定义一个 `<bean>` 标签，就相当于构建了一个 `BeanDefinition`

##### 3.1.1 编写bean与xml

~~~java
public class Person {
    
    private String name;
    
    public String getName() {
        return name;
    }
    
    public void setName(String name) {
        this.name = name;
    }
}
~~~



~~~xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd">

    <bean id="person" class="com.linkedbear.spring.definition.a_quickstart.bean.Person">
        <property name="name" value="zhangsan"/>
    </bean>
</beans>
~~~

##### 3.1.2 获得BeanDefinition

发现`ClassPathXmlApplicationContext`，并没有获得`BeanDefinition`的方法。

![image-20230424210405979](https://wuxie-image.oss-cn-chengdu.aliyuncs.com/2023/04/14/image-20230424210405979.png)

那么我们前面得知`Application` 最终是组合了一个`BeanFactory`。查看`DefaultListableBeanFactory` 得到了该方法

~~~java
public BeanDefinition getBeanDefinition(String beanName) throws NoSuchBeanDefinitionException {
		BeanDefinition bd = this.beanDefinitionMap.get(beanName);
		if (bd == null) {
			if (logger.isTraceEnabled()) {
				logger.trace("No bean named '" + beanName + "' found in " + this);
			}
			throw new NoSuchBeanDefinitionException(beanName);
		}
		return bd;
	}
~~~



所以主类中使用beanfactory获得beanDefinition

~~~java
ClassPathXmlApplicationContext ctx = new ClassPathXmlApplicationContext("BeanDefinition.xml");
BeanDefinition person = ctx.getBeanFactory().getBeanDefinition("person");
System.out.println(person);
/*
Generic bean: class [org.example.BeanDefinition.Person]; scope=; abstract=false; lazyInit=false; autowireMode=0; dependencyCheck=0; autowireCandidate=true; primary=false; factoryBeanName=null; factoryMethodName=null; initMethodName=null; destroyMethodName=null; defined in class path resource [BeanDefinition.xml]
*/

System.out.println(person instanceof GenericBeanDefinition); // true
~~~

#### 3.2 基于@Component

 `Person` 类上添加`@Component` 注解

修改启动类

（注意 `AnnotationConfigApplicationContext` 可以直接调用 `getBeanDefinition` 方法哦，因为继承了`GenericApplicationContext`）

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext("org.example.BeanDefinition");

BeanDefinition person = ctx.getBeanDefinition("person");
System.out.println(person);
System.out.println(person.getClass().getName());

/*
Generic bean: class [org.example.BeanDefinition.Person]; scope=singleton; abstract=false; lazyInit=null; autowireMode=0; dependencyCheck=0; autowireCandidate=true; primary=false; factoryBeanName=null; factoryMethodName=null; initMethodName=null; destroyMethodName=null; defined in file [E:\LearnNote\Spring-Depth-Learning\spring-demo1\target\classes\org\example\BeanDefinition\Person.class]
org.springframework.context.annotation.ScannedGenericBeanDefinition
*/
~~~



> 可以发现，`BeanDefinition` 的打印信息里，最大的不同是加载来源：**基于 xml 解析出来的 bean ，定义来源是 xml 配置文件；基于 `@Component` 注解解析出来的 bean ，定义来源是类的 .class 文件中。**



#### 3.3 基于@Bean

配置类

~~~java
@Configuration
public class BeanDefinitionQuickstartConfiguration {
    
    @Bean
    public Person person() {
        return new Person();
    }
}
~~~



修改主类

~~~java
AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(BeanDefinitionQuickstartConfiguration.class)
BeanDefinition person = ctx.getBeanFactory().getBeanDefinition("person");
System.out.println(person.getClass().getName());
System.out.println(person );
/*
org.springframework.context.annotation.ConfigurationClassBeanDefinitionReader$ConfigurationClassBeanDefinition
Root bean: class [null]; scope=; abstract=false; lazyInit=null; autowireMode=3; dependencyCheck=0; autowireCandidate=true; primary=false; factoryBeanName=beanDefinitionQuickstartConfiguration; factoryMethodName=person; initMethodName=null; destroyMethodName=(inferred); defined in org.example.BeanDefinition.BeanDefinitionQuickstartConfiguration
*/
~~~

具体区别可以发现有这么几个：

- Bean 的类型是 Root bean （ `ConfigurationClassBeanDefinition` 继承自 `RootBeanDefinition` ）
- Bean 的 className 不见了
- 自动注入模式为 `AUTOWIRE_CONSTRUCTOR` （构造器自动注入）
  - `int AUTOWIRE_CONSTRUCTOR = 3;`
- 有 factoryBean 了：person 由 `beanDefinitionQuickstartConfiguration` 的 `person` 方法创建





#### 3.4  BeanDefinition生成原理

1. 通过 xml 加载的 `BeanDefinition` ，它的读取工具是 `XmlBeanDefinitionReader` ，它会解析 xml 配置文件，最终来到 `DefaultBeanDefinitionDocumentReader` 的 `doRegisterBeanDefinitions` 方法，根据 xml 配置文件中的 bean 定义构造 `BeanDefinition` ，最底层创建 `BeanDefinition` 的位置在 `org.springframework.beans.factory.support.BeanDefinitionReaderUtils#createBeanDefinition` 。
2. 通过模式注解 + 组件扫描的方式构造的 `BeanDefinition` ，它的扫描工具是 `ClassPathBeanDefinitionScanner` ，它会扫描指定包路径下包含特定模式注解的类，核心工作方法是 `doScan` 方法，它会调用到父类 `ClassPathScanningCandidateComponentProvider` 的 `findCandidateComponents` 方法，创建 `ScannedGenericBeanDefinition` 并返回。
3. 通过配置类 + `@Bean` 注解的方式构造的 `BeanDefinition` 最复杂，它涉及到配置类的解析。配置类的解析要追踪到 `ConfigurationClassPostProcessor` 的 `processConfigBeanDefinitions` 方法，它会处理配置类，并交给 `ConfigurationClassParser` 来解析配置类，取出所有标注了 `@Bean` 的方法。随后，这些方法又被 `ConfigurationClassBeanDefinitionReader` 解析，最终在底层创建 `ConfigurationClassBeanDefinition` 并返回。



### 4.总结

1. 如何理解 BeanDefinition ？ 
2. BeanDefinition 中主要包含哪些信息？ 
3. BeanDefinition 有哪些类型？分别都有什么特征？

> BeanDefinition 是 Spring 框架中的一个核心概念，它是用来描述 Bean 的配置元数据的对象。BeanDefinition 中包含了 Bean 对象创建所需的所有信息，例如类名、构造函数参数、属性值等。
>
> 具体来说，BeanDefinition 主要包含以下信息：
>
> 1. Bean 类的全限定名或 Class 对象
> 2. 是否是抽象类
> 3. 是否为单例模式
> 4. 是否需要延迟初始化
> 5. 构造函数及参数
> 6. 属性值和依赖项
> 7. 生命周期回调方法等
>
> 在 Spring 中，BeanDefinition 有多种类型，例如 GenericBeanDefinition、RootBeanDefinition、ChildBeanDefinition 等。其中，GenericBeanDefinition 是最常用的类型，它是一个基本的 BeanDefinition 实现，可以描述任何类型的 Bean。而 RootBeanDefinition 和 ChildBeanDefinition 则是继承自 GenericBeanDefinition 的特殊类型，用于表示父子关系的 Bean。
>
> 除此之外，在 Spring 中还有一些其他的 BeanDefinition 实现，例如 AnnotatedBeanDefinition、ScannedGenericBeanDefinition 等，它们都是为了实现不同的 Bean 注入方式而设计的。在使用 Spring 的时候，开发者可以根据需要选择合适的 BeanDefinition 类型，并设置相应的属性来配置 Bean 的创建过程。



## 18. BeanDefinitionRegistry概述【理解】

### 1.介绍

由于官方文档中并没有提及 `BeanDefinitionRegistry` 的设计，故我们只尝试从 javadoc 中获取一些信息。

> 包含 bean 定义的注册表的接口（例如 `RootBeanDefinition` 和 `ChildBeanDefinition` 实例）。通常由内部与 `AbstractBeanDefinition` 层次结构一起工作的 `BeanFactorty` 实现。 这是 SpringFramework 的 bean 工厂包中唯一封装了 bean 的定义注册的接口。标准 `BeanFactory` 接口仅涵盖对完全配置的工厂实例的访问。 `BeanDefinition` 的解析器希望可以使用此接口的实现类来支撑逻辑处理。SpringFramework 中的已知实现者是 `DefaultListableBeanFactory` 和 `GenericApplicationContext` 。

#### 1.1 BeanDefinitionRegistry中存放了所有BeanDefinition

Registry 有注册表的意思，联想下 Windows 的注册表，它存放了 Windows 系统中的应用和设置信息。如果按照这个设计理解，那 `BeanDefinitionRegistry` 中存放的就应该是 `BeanDefinition` 的设置信息。其实 SpringFramework 中的底层，对于 `BeanDefinition` 的注册表的设计，就是一个 **`Map`** ：

```java
// 源自DefaultListableBeanFactory
private final Map<String, BeanDefinition> beanDefinitionMap = new ConcurrentHashMap<>(256);
```

#### 1.2 BeanDefinitionRegistry中维护了BeanDefinition

另外，Registry 还有注册器的意思，既然 Map 有增删改查，那作为 `BeanDefinition` 的注册器，自然也会有 `BeanDefinition` 的注册功能咯。`BeanDefinitionRegistry` 中有 3 个方法，刚好对应了 `BeanDefinition` 的增、删、查：

```java
void registerBeanDefinition(String beanName, BeanDefinition beanDefinition)
            throws BeanDefinitionStoreException;

void removeBeanDefinition(String beanName) throws NoSuchBeanDefinitionException;

BeanDefinition getBeanDefinition(String beanName) throws NoSuchBeanDefinitionException;
```

#### 1.3 BeanDefinitionRegistry支撑其它组件运行

> `BeanDefinition` 的加载器希望可以使用此接口的实现类来支撑逻辑处理。

javadoc 中的 Reader 可以参照上一章提到了 `XmlBeanDefinitionReader` ，它是用来读取和加载 xml 配置文件的组件。加载 xml 配置文件的目的就是读取里面的配置，和定义好要注册到 IOC 容器的 bean 。`XmlBeanDefinitionReader` 要在加载完 xml 配置文件后，将配置文件的流对象也好，文档对象也好，交给解析器来解析 xml 文件，解析器拿到 xml 文件后要解析其中定义的 bean ，并且封装为 `BeanDefinition` 注册到 IOC 容器，这个时候就需要 `BeanDefinitionRegistry` 了。所以在这个过程中，**`BeanDefinitionRegistry` 会支撑 `XmlBeanDefinitionReader` 完成它的工作**。

当然，`BeanDefinitionRegistry` 不止支撑了这一个哈，还记得之前小册 17 章，学习模块装配时用到的 `ImportBeanDefinitionRegistrar` 吗？它的 `registerBeanDefinitions` 方法是不是也传入了一个 `BeanDefinitionRegistry` 呀？所以说这个 `BeanDefinitionRegistry` 用到的位置还是不少的，小伙伴们要予以重视哦。

#### 1.4 BeanDefinitionRegistry的主要实现是DefaultListableBeanFactory

注意这个地方我没说是唯一实现哦，是因为 `BeanDefinitionRegistry` 除了有最最常用的 `DefaultListableBeanFactory` 之外，还有一个不常用的 `SimpleBeanDefinitionRegistry` ，但这个 `SimpleBeanDefinitionRegistry` 基本不会去提它，是因为这个设计连内部的 IOC 容器都没有，仅仅是一个 `BeanDefinitionRegistry` 的表面实现而已，所以我们当然不会用它咯。

可能有的小伙伴借助 IDE 发现很多 `ApplicationContext` 也实现了它，但我想请这部分小伙伴回想一下，`ApplicationContext` 本身管理 Bean 吗？不吧，`ApplicationContext` 不都是内部组合了一个 `DefaultListableBeanFactory` 来实现的嘛，所以我们说，唯一真正落地实现的是 `DefaultListableBeanFactory` 这话是正确合理的。

#### 1.5 【面试题】面试中如何概述BeanDefinitionRegistry

以下答案仅供参考，可根据自己的理解调整回答内容：

**`BeanDefinitionRegistry` 是维护 `BeanDefinition` 的注册中心，它内部存放了 IOC 容器中 bean 的定义信息，同时 `BeanDefinitionRegistry` 也是支撑其它组件和动态注册 Bean 的重要组件。在 SpringFramework 中，`BeanDefinitionRegistry` 的实现是 `DefaultListableBeanFactory` 。**

### 2. BeanDefinitionRegistry维护BeanDefinition的使用【熟悉】

对于 `BeanDefinitionRegistry` 内部的设计，倒是没什么好说的，主要还是研究它如何去维护 `BeanDefinition` 。

> 本小节源码位置：`com.linkedbear.spring.definition.b_registry`

#### 2.1 BeanDefinition的注册

对于 `BeanDefinition` 的注册，目前我们接触到的方式是在 17 章模块装配中使用的 `ImportBeanDefinitionRegistrar` ：

```java
public class WaiterRegistrar implements ImportBeanDefinitionRegistrar {
    
    @Override
    public void registerBeanDefinitions(AnnotationMetadata metadata, BeanDefinitionRegistry registry) {
        registry.registerBeanDefinition("waiter", new RootBeanDefinition(Waiter.class));
    }
}
```

之前的这个例子中是直接 **new** 了一个 `RootBeanDefinition` ，其实 `BeanDefinition` 的构造可以借助**建造器**生成，下面我们再演示一个例子。

#### 2.1.1 声明Person类

像往常一样，搞一个比较简单的 `Person` 就好啦，记得声明几个属性和 `toString` 方法：

```java
public class Person {
    
    private String name;
    
    public String getName() {
        return name;
    }
    
    public void setName(String name) {
        this.name = name;
    }
    
    @Override
    public String toString() {
        return "Person{" + "name='" + name + '\'' + '}';
    }
}
```

#### 2.1.2 编写ImportBeanDefinitionRegistrar的实现类

编写一个 `PersonRegister` ，让它实现 `ImportBeanDefinitionRegistrar` ，这样就可以拿到 `BeanDefinitionRegistry` 了：

```java
public class PersonRegister implements ImportBeanDefinitionRegistrar {
    
    @Override
    public void registerBeanDefinitions(AnnotationMetadata importingClassMetadata, BeanDefinitionRegistry registry) {
        registry.registerBeanDefinition("person",
                BeanDefinitionBuilder.genericBeanDefinition(Person.class).addPropertyValue("name", "zhangsan")
                        .getBeanDefinition());
    }
}
```

注意这里面的写法，使用 `BeanDefinitionBuilder` ，是可以创建 `GenericBeanDefinition` 、`RootBeanDefinition` 和 `ChildBeanDefinition` 三种类型的，此处小册使用 `GenericBeanDefinition` ，后续直接向 `BeanDefinition` 中添加 bean 中属性的值就好，整个构造过程一气呵成，非常的简单。

#### 2.1.3 编写配置类导入PersonRegister

编写一个配置类，把上面刚写好的 `PersonRegister` 导入进去（这一步不要忘了哦）：

```java
@Configuration
@Import(PersonRegister.class)
public class BeanDefinitionRegistryConfiguration {
    
}
```

#### 2.1.4 测试获取Person

万事俱备，下面编写测试启动类，使用 `BeanDefinitionRegistryConfiguration` 驱动 IOC 容器，并从容器中取出 `Person` 并打印：

```java
public class BeanDefinitionRegistryApplication {
    
    public static void main(String[] args) throws Exception {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext(
                BeanDefinitionRegistryConfiguration.class);
        Person person = ctx.getBean(Person.class);
        System.out.println(person);
    }
}
```

运行 `main` 方法，控制台中打印了 `Person` 的 name 属性是有值的，说明 SpringFramework 已经按照我们预先定义好的 `BeanDefinition` ，注册到 IOC 容器，并且生成了对应的 Bean 。

```ini
Person{name='zhangsan'}
```

#### 2.2 BeanDefinition的移除

`BeanDefinitionRegistry` 除了能给 IOC 容器中添加 `BeanDefinition` ，还可以移除掉一些特定的 `BeanDefinition` 。这种操作可以在 Bean 的实例化之前去除，以阻止 IOC 容器创建。

要演示 `BeanDefinition` 的移除，需要一个现阶段没见过的 API ，咱们先学着用一下，到后面我们会系统的学习它的用法。

> 本小节源码位置：`com.linkedbear.spring.definition.c_removedefinition`

##### 2.2.1 声明Person

这次声明的 `Person` 类要加一个特殊的属性：**sex** ，性别，它在后面会起到判断作用。

声明好 getter 、setter 和 `toString` 方法即可。

```java
public class Person {
    
    private String name;
    private String sex;
    
    // getter 、setter 、 toString
}
```

##### 2.2.2 声明配置类

接下来要注册两个 `Person` ，分别注册一男一女。

由上一章 `BeanDefinition` 的注册方式与实现类型，可知如果此处使用注解配置类的方式注册 Bean ( `@Bean` ) ，生成的 `BeanDefinition` 将无法取到 `beanClassName` （也无法取到 PropertyValues ），故此处选用 xml 方式注册 Bean 。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xmlns:context="http://www.springframework.org/schema/context"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd http://www.springframework.org/schema/context https://www.springframework.org/schema/context/spring-context.xsd">

    <bean id="aqiang" class="com.linkedbear.spring.definition.c_removedefinition.bean.Person">
        <property name="name" value="阿强"/>
        <property name="sex" value="male"/>
    </bean>

    <bean id="azhen" class="com.linkedbear.spring.definition.c_removedefinition.bean.Person">
        <property name="name" value="阿珍"/>
        <property name="sex" value="female"/>
    </bean>

    <!-- 注意此处要开启包扫描 -->
    <context:component-scan base-package="com.linkedbear.spring.definition.c_removedefinition.config"/>
</beans>
```

##### 2.2.3 编写剔除BeanDefinition的后置处理器

这里涉及到后置处理器的概念了，没见过没关系，不会搞没关系，先照着葫芦画瓢，后面马上就学到了。

要剔除 `BeanDefinition` ，需要实现 `BeanFactoryPostProcessor` 接口，并重写 `postProcessBeanFactory` 方法：（记得标注 `@Component` 注解哦）

```java
@Component
public class RemoveBeanDefinitionPostProcessor implements BeanFactoryPostProcessor {
    
    @Override
    public void postProcessBeanFactory(ConfigurableListableBeanFactory beanFactory) throws BeansException {
    
    }
}
```

注意方法的入参，它是一个 `ConfigurableListableBeanFactory` ，不用想，它的唯一实现一定是 `DefaultListableBeanFactory` 。又从前面了解到 `DefaultListableBeanFactory` 实现了 `BeanDefinitionRegistry` 接口，所以这里我们就可以直接将 `beanFactory` 强转为 `BeanDefinitionRegistry` 类型。

于是，我们就可以编写如下的剔除逻辑：**移除 IOC 容器中所有性别为 male 的 Person** 。

```java
@Override
public void postProcessBeanFactory(ConfigurableListableBeanFactory beanFactory) throws BeansException {
    BeanDefinitionRegistry registry = (BeanDefinitionRegistry) beanFactory;
    // 获取IOC容器中的所有BeanDefinition
    for (String beanDefinitionName : beanFactory.getBeanDefinitionNames()) {
        // 判断BeanDefinition对应的Bean是否为Person类型
        BeanDefinition beanDefinition = beanFactory.getBeanDefinition(beanDefinitionName);
        if (Person.class.getName().equals(beanDefinition.getBeanClassName())) {
            // 判断Person的性别是否为male
            // 使用xml配置文件对bean进行属性注入，最终取到的类型为TypedStringValue，这一点不需要记住
            TypedStringValue sex = (TypedStringValue) beanDefinition.getPropertyValues().get("sex");
            if ("male".equals(sex.getValue())) {
                // 移除BeanDefinition
                registry.removeBeanDefinition(beanDefinitionName);
            }
        }
    }
}
```

##### 2.2.4 测试获取“阿强”

这一次我们又要用 `ClassPathXmlApplicationContext` 来加载配置文件驱动 IOC 容器了，写法很简单，直接从 IOC 容器中取 “aqiang” 就好：

```java
public class RemoveBeanDefinitionApplication {
    
    public static void main(String[] args) throws Exception {
        ClassPathXmlApplicationContext ctx = new ClassPathXmlApplicationContext("definition/remove-definitions.xml");
        Person aqiang = (Person) ctx.getBean("aqiang");
        System.out.println(aqiang);
    }
}
```

运行 `main` 方法，控制台打印 `NoSuchBeanDefinitionException` 的异常，证明 “aqiang” 对应的 `BeanDefinition` 已经被移除了，无法创建 `Person` 实例。

好了，到这里，对 `BeanDefinitionRegistry` 有一个比较清晰的认识就好，具体操作不需要太深入了解，会用就够啦。

### 3. BeanDefinition的合并【了解】

了解完 `BeanDefinitionRegistry` ，回过头来再学习一个 `BeanDefinition` 的特性：**合并**。

关于合并这个概念，可能有些小伙伴没有概念，小册先来解释一下合并的意思。

#### 3.1 如何理解BeanDefinition的合并

上一章我们知道，之前在 xml 配置文件中定义的那些 bean ，最终都转换为一个个的 `GenericBeanDefinition` ，它们都是相互独立的。比如这样：

```xml
<bean class="com.linkedbear.spring.basic_dl.b_bytype.bean.Person"></bean>
<bean class="com.linkedbear.spring.basic_dl.b_bytype.dao.impl.DemoDaoImpl"/>
```

但其实，bean 也是存在**父子关系**的。与 Class 的抽象、继承一样，`<bean>` 标签中有 **abstract** 属性，有 **parent** 属性，由此就可以形成父子关系的 `BeanDefinition` 了。

下面小册演示一个实例，讲解 `BeanDefinition` 的合并。

#### 3.2 BeanDefinition合并的体现

先构建一个比较简单的场景吧：所有的**动物**都归**人**养，动物分很多种（猫啊 狗啊 猪啊 巴拉巴拉）。

下面我们基于这个场景来编码演绎。

##### 3.2.1 声明实体类

对于这几个实体类，前面已经写过很多次了，这里快速编写出来就 OK ：

```java
public class Person {
    
}
public abstract class Animal {

    private Person person;
    
    public Person getPerson() {
        return person;
    }
    
    public void setPerson(Person person) {
        this.person = person;
    }
}
```

`Cat` 要继承自 `Animal` ，并且为了方便打印出 `person` ，这里就不直接使用 IDEA 的 `toString` 方法生成了，而是在此基础上改造一下：

```java
public class Cat extends Animal {
    
    private String name;
    
    public String getName() {
        return name;
    }
    
    public void setName(String name) {
        this.name = name;
    }
    
    @Override
    public String toString() {
        return "Cat{" + "name=" + name + ", person='" + getPerson() + '\'' + "}";
    }
}
```

##### 3.2.2 编写xml配置文件

要体现 `BeanDefinition` 的合并，要使用配置文件的形式，前面也说过了。那下面咱就来造一个配置文件，先把 `Person` 注册上去。

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd">

    <bean id="person" class="com.linkedbear.spring.definition.d_merge.bean.Person"/>
</beans>
```

接下来要注册 `Animal` 和 `Cat` 了。按照之前的写法，这里只需要注册 `Cat` 就可以了，像这样写就 OK ：

```xml
<bean class="com.linkedbear.spring.definition.d_merge.bean.Cat" parent="abstract-animal">
    <property name="person" ref="person"/>
    <property name="name" value="咪咪"/>
</bean>
```

但试想，如果要创建的猫猫狗狗猪猪太多的话，每个 bean 都要注入 property ，这样可不是好办法。由此，就可以使用 `BeanDefinition` 合并的特性来优化这个问题。

我们直接在 xml 中注册一个 `Animal` ：

```xml
<bean class="com.linkedbear.spring.definition.d_merge.bean.Animal"></bean>
```

但这样写完之后，IDEA 会报红，给出这样的提示：

![img](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/1811045a723f455baf9e49a8701fb0b1~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

很明显嘛，抽象类怎么能靠一个 `<bean>` 标签构造出对象呢？所以，`<bean>` 标签里有一个属性，就是标注这个 bean 是否是抽象类：

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/8296e2db237648098d16e8411c6230a2~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

如此，咱就可以把这个 `Animal` 声明好了，由于是 **abstract** 类型的 bean ，那也就可以搞定注入的事了：

```xml
<bean id="abstract-animal" class="com.linkedbear.spring.definition.d_merge.bean.Animal" abstract="true">
    <property name="person" ref="person"/>
</bean>
```

接下来要声明 `Cat` 了，有 **abstract** 就有 **parent** ，想必不用我多说小伙伴们也能猜到如何写了：

```xml
<bean id="cat" class="com.linkedbear.spring.definition.d_merge.bean.Cat" parent="abstract-animal">
    <property name="name" value="咪咪"/>
</bean>
```

这里就不再需要声明 `person` 属性的注入了，因为继承了 `abstract-animal` ，相应的依赖注入也就都可以继承过来。

这样 xml 配置文件就写完了。

##### 3.2.3 测试运行

编写启动类，使用 xml 配置文件驱动 IOC 容器，并从 `BeanFactory` 中取出 cat 的 `BeanDefinition` ：

```java
public class MergeBeanDefinitionApplication {
    
    public static void main(String[] args) throws Exception {
        ClassPathXmlApplicationContext ctx = new ClassPathXmlApplicationContext("definition/definition-merge.xml");
        Cat cat = (Cat) ctx.getBean("cat");
        System.out.println(cat);
        
        BeanDefinition catDefinition = ctx.getBeanFactory().getBeanDefinition("cat");
        System.out.println(catDefinition);
    }
}
```

运行 `main` 方法，发现 `Cat` 里确实注入了 `person` 对象，可是获取出来的 `BeanDefinition` ，除了有了一个 `parentName` 之外，跟普通的 bean 没有任何不一样的地方。

```ini
Cat{name=咪咪, person='com.linkedbear.spring.definition.d_merge.bean.Person@31dc339b'}
Generic bean with parent 'abstract-animal': class [com.linkedbear.spring.definition.d_merge.bean.Cat]; scope=;   ......(太长省略)
```

可能会有小伙伴产生疑惑了：这就算是 `BeanDefinition` 的合并了吗？哪里有体现呢？要么我 Debug 看下结构？

以 Debug 的形式重新运行 `main` 方法，发现获取到的 `catDefinition` 里并没有把 `person` 的依赖带进来：

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/173fbedcb5674458a6f6a21a33ff8357~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

哦，合着并没有合并 `BeanDefinition` 呗？那这一套花里胡哨的搞蛇皮呢？

等一下，先冷静冷静，会不会是我们的方法不对呢？既然是 `BeanDefinition` 的合并，那不加个 **merge** 的关键字，好意思说是合并吗？

试着重新调一下方法，发现 `ConfigurableListableBeanFactory` 里竟然也有一个 `getMergedBeanDefinition` 方法！它来自 `ConfigurableBeanFactory` ，它就是用来**将本身定义的 bean 定义信息，与继承的 bean 定义信息进行合并后返回**的。

##### 3.2.4 换用getMergedBeanDefinition

修改下测试运行，将 `getBeanDefinition` 换为 `getMergedBeanDefinition` ，重新运行 `main` 方法，发现控制台打印的 `BeanDefinition` 的类型变为了 `RootBeanDefinition` ，而且也没有 `parentName` 相关的信息了：

```java
Root bean: class [com.linkedbear.spring.definition.d_merge.bean.Cat]; scope=singleton; abstract=false; lazyInit=false; autowireMode=0; dependencyCheck=0; autowireCandidate=true; primary=false; factoryBeanName=null; factoryMethodName=null; initMethodName=null; destroyMethodName=null; defined in class path resource [definition/definition-merge.xml]
```

以 Debug 方式运行，此时的 `propertyvalues` 中已经有两个属性键值对了：

![img](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/be105a8a0d7f4c7c906783887d990d90~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

### 4. 设计BeanDefinition的意义【理解】

看到这里，小伙伴们可能会有一个大大的问号，也有可能是大大的感叹号，那就是：**SpringFramework 为什么会设计 `BeanDefinition` 呢？直接注册 Bean 不好吗？**理解这样的一个设计：**定义信息 → 实例**。

像我们平时编写 **Class** 再 **new** 出对象一样，**SpringFramework 面对一个应用程序，它也需要对其中的 bean 进行定义抽取，只有抽取成可以统一类型 / 格式的模型，才能在后续的 bean 对象管理时，进行统一管理，也或者是对特定的 bean 进行特殊化的处理。而这一切的一切，最终落地到统一类型上，就是 `BeanDefinition` 这个抽象化的模型。**
