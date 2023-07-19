import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.8.21"
    application
}

group = "social.miauw"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(kotlin("test"))
    implementation("com.rabbitmq:amqp-client:5.18.0")
    implementation("org.slf4j:slf4j-simple:2.0.7")
    implementation("org.slf4j:slf4j-api:2.0.7")
    implementation("com.google.code.gson:gson:2.10.1")
    implementation("org.postgresql:postgresql:42.6.0")
    implementation("at.favre.lib:bcrypt:0.10.2")
    implementation("org.bitbucket.b_c:jose4j:0.9.3")
    implementation("org.json:json:20230618")
    implementation("org.jdbi:jdbi3-kotlin:3.40.0-a0")
    implementation("org.jdbi:jdbi3-postgres:3.40.0-a0")
    implementation("org.jdbi:jdbi3-core:3.40.0-a0")
}

tasks.test {
    useJUnitPlatform()
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "11"
}

application {
    mainClass.set("MainKt")
}