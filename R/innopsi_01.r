setwd("/Projects/innopsi/R")
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





















mode(data.testing$x2)

# integer values 0,1,2
# numerica values 0 to 100

hist(data.testing$x2)


"
One Approach - data extensive, brute force
  
create subsets of the data, and possible recursive subsets
for example, for the integer values, 6 subsets each, for the numeric values 5 subsets

with 20 integer
and  20 numeric values

single column subsets count 20*6*20*5 = 1800 subsets
now taking all two column subsets 1800^2
and so on 
this grows exponentially quickly

or starting with initial 1800 subsets  
branch and test, .. same complexity
pick x1,  test, pick x1 and x2, pick x1 and x3
it is simple counting though at this stage

for group and its complement
average fitness < average fitness + .6

if group is found,
persist group and fitness and continue (or stop and move onto next subset)

might be better done with C#

"

# select each column, and it's six subsets and evaluate the difference in the values
d2 <- data.testing[data.testing$dataset==2,]

d2.0 <- d2[d2$trt==0 & d2$x1<1,c(1:4,5)]
d2.1 <- d2[d2$trt==1 & d2$x1<1,c(1:4,5)]
m.0 <- mean(d2.0$y)
m.1 <- mean(d2.1$y)
m.0 - m.1

dt.0 <-data.testing[data.testing$dataset==2 & data.testing$x1 ==0,]
dt.1 <-data.testing[data.testing$dataset==2 & data.testing$x1 ==1,]
m.0 <- mean(dt.0$y)
m.1 <- mean(dt.1$y)
m.0-m.1

dt.0 <-data.testing[data.testing$dataset==2 & data.testing[[5]] ==0,]
dt.1 <-data.testing[data.testing$dataset==2 & data.testing[[5]] ==1,]
m.0 <- mean(dt.0$y)
m.1 <- mean(dt.1$y)
m.0-m.1



# check each column individually
i <- c(5:24)
lcc <-c()
for(v in i){
  # one subset
  dt.0.1 <-data.testing[data.testing$dataset==2 & data.testing$trt==0 & data.testing[[v]] ==0,]
  dt.1.1 <-data.testing[data.testing$dataset==2 & data.testing$trt==1 & data.testing[[v]] ==0,]
  
  dt.0.2 <-data.testing[data.testing$dataset==2 & data.testing$trt==0 & data.testing[[v]] ==1,]
  dt.1.2 <-data.testing[data.testing$dataset==2 & data.testing$trt==1 & data.testing[[v]] ==1,]
  
  dt.0.3 <-data.testing[data.testing$dataset==2 & data.testing$trt==0 & data.testing[[v]] ==2,]
  dt.1.3 <-data.testing[data.testing$dataset==2 & data.testing$trt==1 & data.testing[[v]] ==2,]
  
  dt.0.4 <-data.testing[data.testing$dataset==2 & data.testing$trt==0 & (data.testing[[v]] ==0  || data.testing[[v]]==1),]
  dt.1.4 <-data.testing[data.testing$dataset==2 & data.testing$trt==1 & (data.testing[[v]] ==0  || data.testing[[v]]==1),]
  
  dt.0.5 <-data.testing[data.testing$dataset==2 & data.testing$trt==0 & (data.testing[[v]] ==0  || data.testing[[v]]==2),]
  dt.1.5 <-data.testing[data.testing$dataset==2 & data.testing$trt==1 & (data.testing[[v]] ==0  || data.testing[[v]]==2),]
  
  
  dt.0.6 <-data.testing[data.testing$dataset==2 & data.testing$trt==0 & (data.testing[[v]] ==1  || data.testing[[v]]==2),]
  dt.1.6 <-data.testing[data.testing$dataset==2 & data.testing$trt==1 & (data.testing[[v]] ==1  || data.testing[[v]]==2),]
  
  
  #m.0.1 <- mean(dt.0$y)
  #m.1.1 <- mean(dt.1$y)
  #m.0-m.1
  
  #print(m.0)
  #print(m.1)
  #print(m.0-m.1)
  lcc[[v]] <- c(
    mean(dt.0.1$y-dt.1.1$y),
    mean(dt.0.2$y-dt.1.2$y),
    mean(dt.0.3$y-dt.1.3$y),
    mean(dt.0.4$y-dt.1.4$y),
    mean(dt.0.5$y-dt.1.5$y),
    mean(dt.0.6$y-dt.1.6$y)
    )
}

t <-data.frame(lcc)















