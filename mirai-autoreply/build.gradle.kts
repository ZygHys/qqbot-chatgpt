plugins {
    val kotlinVersion = "1.5.30"
    kotlin("jvm") version kotlinVersion
    kotlin("plugin.serialization") version kotlinVersion

    id("net.mamoe.mirai-console") version "2.11.1"
}

group = "per.autumn.mirai.autoreply"
version = "1.3.2"

repositories {
    maven("https://maven.aliyun.com/repository/public")
    mavenCentral()
}
dependencies {
    implementation("cn.hutool:hutool-all:5.8.0")
    implementation("net.objecthunter:exp4j:0.4.8")
    implementation ("com.google.code.gson:gson:2.8.5")
    testImplementation(kotlin("test"))
}
tasks.test {
    useJUnitPlatform()
}

