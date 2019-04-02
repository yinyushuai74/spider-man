# spider-man

##robot.txt:

*https://www.wongnai.com/robots.txt

*https://www.vietnammm.com/robots.txt

*https://www.now.vn/robots.txt


========================================
#intro
spider-man was developed to scrape data from now.vn, this is a guide to help u run it in your local.

#Step 1:Download the spider-man
for mac :spider-man

for windows:spider-man.exe

#step 2:
Open a terminal, and go to the directory where the spider-man is downloaded

cd ~/Downloads



#step 3:run the spider-man
./spider-man -size={{size}} -city={{cityID}} -mode={{mode}} -district={{districtID}} -sortType={{sortType}}

e.g. scrape the data of Ha Noi City, get the cityID 218 from above

 ./spider-man -city=218
arguments in command:
-size : How many rows for each csv file default:5000 rows per file
-city : Which city is your target city default: HCM(217)
 HCM City 217

Ha Noi City 218

Da Nang City 219

Can Tho City 221

Hai Phong City 220

Hue City 273

Khanh Hoa 248

Dong Nai 222

Nghe An 257

Vung Tau 223

Binh Duong 230

Lam Dong 254

Quang Ninh 265

Quang Nam 263

-mode : scrape mode default simple mode
  0: simple mode would only scrape merchant info ( no menu)

  1: elaborate mode would scrape menu (merchant info and menu)

-sortType:sort type 
recent(normal, default):10

verified: 26

opening: 6

shiped by now: 7

featured:1

preferred:25

VAT:28

top:2

new:29

-district: only scrape the target district (default: do not filter specific district scrape all merchant of the city )
Note ：district must belong with city，if not will stop.

e.g.  if run ./spider-man -size=20 -city=218 -mode=0 -district=5 -sortType=26  districtID= 5 for Quận 3 should  belong HCM(217)  so if -city=218 will return no merchant.         

district list:


#popular command:
./spider-man -city=217 -mode=1  

all the merchant in HCM and get the merchant with menu

./spider-man -city=217  -sortType=26

all the verified merchant in HCM 

./spider-man -city=217  -size=10000 -mode=1

 get the all the merchant with menu in HCM and u want the data separate of 10000 rows each file



#Note:
the merchant information would display in your terminate and wait it, when it finished you can see `put merchantID done` in your terminate and you can use `control+c` to stop it.





then the file name of merchant_{{cityname}}_{{index}}.csv would be display( if interrupt when the scrape unfinished the file would be blank).










