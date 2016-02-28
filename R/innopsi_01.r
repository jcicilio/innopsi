setwd("/Projects/GoLang/src/innopsi/R")
fname = "innopsi_01.csv"
saveName = paste0("/Projects/innopsi/data/ModelResults/", fname)

# load machine settings
data.testing = read.csv("../data/InnoCentive_9933623_Training_Data.csv", header = TRUE)
min(data.testing)
max(data.testing)
str(data.testing)
summary(data.testing)


d1 <- subset(data.testing, dataset==1)
summary(d1)

# Take all subsets of x1 through 20
# create a dataset of all subsets 6 levels 20 rows 120 subsets, and evaluate them
# data format  dataset, subsetid, column number, value, y
d <- data.frame()

# loop through rows, and then row bind to datasetS
# or 5 columns to start
dv.1 <-subset(data.testing, dataset==1)
dv.1 <- dv.1[,c(1,2,3,4,5)]
dv.1.t <-subset(dv.1, trt==1)
dv.1.c <-subset(dv.1, trt==0)
summary(dv.1.t)
summary(dv.1.c)


# loop through rows, and then row bind to datasetS
# or 5 columns to start
dv.1 <-subset(data.testing, dataset==1)
dv.1 <- dv.1[,c(1,2,3,4,5)]
dv.1.t <-subset(dv.1, trt==1)
dv.1.c <-subset(dv.1, trt==0)
summary(dv.1.t)
summary(dv.1.c)
mean(dv.1.t$y)/sd(dv.1.t$y)
mean(dv.1.c$y)/sd(dv.1.t$y)


dv.2 <-subset(data.testing, dataset==1)
dv.2 <- dv.2[,c(1,2,3,4,6)]
dv.2.t <-subset(dv.2, trt==1)
#dv.2.t <- dv.2.t[dv.2.t$x2==1 | dv.2.t$x2==2, ]
dv.2.c <-subset(dv.2, trt==0)
summary(dv.2.t)
summary(dv.2.c)
mean(dv.2.t$y)/sd(dv.2.t$y)
mean(dv.2.c$y)/sd(dv.2.t$y)


dv.3 <-subset(data.testing, dataset==1 )
dv.3 <- dv.3[,c(1,2,3,4,9)]
dv.3.t <-subset(dv.3, trt==1)
dv.3.t <- dv.3.t[dv.3.t$x5==1 |dv.3.t$x5==2, ]
dv.3.c <-subset(dv.3, trt==0)
dv.3.c <- dv.3.c[dv.3.c$x5==1 |dv.3.c$x5==2, ]
summary(dv.3.t)
summary(dv.3.c)
mtsd<-mean(dv.3.t$y)/sd(dv.3.t$y)
mcsd<-mean(dv.3.c$y)/sd(dv.3.t$y)
sd(dv.3.t$y)
sd(dv.3.c$y)
(mtsd-mcsd)/sd(dv.3.t$y)


dv.2 <-subset(data.testing, dataset==3)
dv.2 <- dv.2[,c(1,2,3,4,6)]
dv.2.t <-subset(dv.2, dv.2$trt==1)
dv.2.c <-subset(dv.2, dv.2$trt==0)
summary(dv.2.t)
summary(dv.2.c)


dv.3 <-subset(data.testing, dataset==3)
dv.3 <- dv.3[,c(1,2,3,4,7)]
dv.3.t <-subset(dv.3, trt==1)
dv.3.c <-subset(dv.3, trt==0)
summary(dv.3.t)
summary(dv.3.c)


# loop through rows, and then row bind to datasetS
# or 5 columns to start
dv.1 <-subset(data.testing, dataset==3)
dv.1 <- dv.1[,c(1,2,3,4,5)]
dv.1.t <-subset(dv.1, trt==1)
dv.1.c <-subset(dv.1, trt==0)
summary(dv.1.t)
summary(dv.1.c)
mean(dv.1.t$y)/sd(dv.1.t$y)
mean(dv.1.c$y)/sd(dv.1.t$y)


dv.2 <-subset(data.testing, dataset==3)
dv.2 <- dv.2[,c(1,2,3,4,6)]
dv.2.t <-subset(dv.2, trt==1)
#dv.2.t <- dv.2.t[dv.2.t$x2==1 | dv.2.t$x2==2, ]
dv.2.c <-subset(dv.2, trt==0)
summary(dv.2.t)
summary(dv.2.c)
mean(dv.2.t$y)/sd(dv.2.t$y)
mean(dv.2.c$y)/sd(dv.2.t$y)


dv.3 <-subset(data.testing, dataset==3 )
dv.3 <- dv.3[,c(1,2,3,4,9)]
dv.3.t <-subset(dv.3, trt==1)
dv.3.t <- dv.3.t[dv.3.t$x5==1 |dv.3.t$x5==2, ]
dv.3.c <-subset(dv.3, trt==0)
dv.3.c <- dv.3.c[dv.3.c$x5==1 |dv.3.t$x5==2, ]
summary(dv.3.t)
summary(dv.3.c)
mtsd<-mean(dv.3.t$y)/sd(dv.3.t$y)
mcsd<-mean(dv.3.c$y)/sd(dv.3.t$y)
sd(dv.3.t$y)
sd(dv.3.t$y)
(mtsd-mcsd)/sd(dv.3.t$y)


# optional score
# (mtsd-mtsd)/sd > .6










# check dataset 3
d3 <- subset(data.testing, dataset==3)
d3.t <- subset(d3, x5==1 | x5==2)

d3.tt <- subset(d3.t, trt==1)
d3.tc <- subset(d3.t, trt==0)

t.test(d3.tt$y)
t.test(d3.tc$y)

mt<-mean(d3.tt$y)
mc<-mean(d3.tc$y)
mean(d3.t0$y)

st<- sd(d3.tt$y)
sc<- sd(d3.tc$y)
sd(d3.t0$y)
st/(st +sc)

summary(d3.tt$y)
summary(d3.tc$y)
summary(d3.t0$y)



k<-(mt-mc)/sd(d3.tt$y)
kp<-(mt-mc)/sd(d3.tc$y)






setwd("/Projects/GoLang/src/innopsi/R")

# load machine settings
data = read.csv("../data/rangeOfScores.csv", header = TRUE)
summary(data)
sd(data$score)
max(data$score) - sd(data$score)
min(data$score) - sd(data$score)
hist(data$score)

(max(data$score) - mean(data$score))/sd(data$score)
(min(data$score) - mean(data$score))/sd(data$score)

# Get this data for test sets 1 and test sets 3 (1 no indicators, 3, x5 is 1 or 2)
# check this value (max-mean)/sd and (min-mean)/sd of scores for all test scores < 0
# goal is to see if there is a pattern to these scores
# use first level 240 


setwd("/Projects/GoLang/src/innopsi/R")

# load machine settings
data = read.csv("../data/testscoreresults240.csv", header = TRUE)
summary(data)

d1 = data[data$dataset==1 & data$score!=0,]
d2 = data[data$dataset==2 & data$score!=0,]
d3 = data[data$dataset==3 & data$score!=0,]
d4 = data[data$dataset==4 & data$score!=0,]
summary(d1$score)
summary(d2$score)
summary(d3$score)
summary(d4$score)
sd(d1$score)
sd(d2$score)
sd(d3$score)
sd(d4$score)
hist(d1$score)
hist(d2$score)
hist(d3$score)
hist(d4$score)
min(d1$score)-max(d1$score)
min(d2$score)-max(d2$score)
min(d3$score)-max(d3$score)
min(d4$score)-max(d4$score)


(min(d1$score)-max(d1$score))/sd(d1$score)
(min(d2$score)-max(d2$score))/sd(d2$score)
(min(d3$score)-max(d3$score))/sd(d3$score)
(min(d4$score)-max(d4$score))/sd(d4$score)

setwd("/Projects/GoLang/src/innopsi/R")

data = read.csv("../data/range20160228.csv", header = TRUE)
summary(data)
sd(data$score)

hist(data$score)
v2 <-2*sd(data$score) -1
v3<-3*sd(data$score) -1
v4 <-4*sd(data$score) -1

dk.2 <- data[data$score > v2,]
dk.3 <- data[data$score > v3,]
dk.4 <- data[data$score > v4,]
# 2 sd  123 of 1200
# 3 sd   56 of 1200


