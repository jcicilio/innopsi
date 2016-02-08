setwd("/Projects/innopsi/R")
fname = "innopsi_01.csv"
saveName = paste0("/Projects/innopsi/data/ModelResults/", fname)

# load machine settings
data.testing = read.csv("../data/InnoCentive_9933623_Data.csv", header = TRUE)
min(data.testing)
max(data.testing)
str(data.testing)
summary(data.testing)

# integer values 0,1,2
# numerica values 0 to 100

hist(data.testing$x33)


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
