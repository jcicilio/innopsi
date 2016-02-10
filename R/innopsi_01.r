setwd("/Projects/innopsi/R")
fname = "innopsi_01.csv"
saveName = paste0("/Projects/innopsi/data/ModelResults/", fname)

# load machine settings
data.testing = read.csv("../data/InnoCentive_9933623_Data.csv", header = TRUE)
min(data.testing)
max(data.testing)
str(data.testing)
summary(data.testing)

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















