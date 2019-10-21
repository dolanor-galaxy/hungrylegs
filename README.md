# HungryLegs 🏊🏻‍🚴‍🏃‍

HungryLegs is your own personal training analysis center.

At the moment HungryLegs is only suitable for software engineers. Or, at least, people comfortable with compiling alpha code, writing raw SQL, and are looking to do their own data analysis.

If you want something non-nerdy, may I suggest:

* [TrainingPeaks](https://www.trainingpeaks.com/) $$
* [Golden Cheetah](https://www.goldencheetah.org/) Free

HungryLegs came about because I wanted some tools to help me with training plans, and also a way for me to analyse data from my runs / bikes / swims.

## Compiling

Running `make` on it's own will give hints on how to build. You're probably after `make build.all.cli` which will make both command line application: `hungrylegs` and `plan`.

## HungryLegs App

The command line application `hungrylegs` will take a directory of `.FIT` or `.TCX` files, and create an sqlite database of them. You can use this database to analyse your activities. `.FIT`  files are typically created by [Garmin](https://www.garmin.com) (amongst other training and fitness products), as are `.TCX` (an XML format similar to FIT).

After you build, edit the `config.json` file. Change the root athlete name, and point the `import` setting to a directory where you have some FIT and / or TCX files are. Then run `hungrylegs`.

(I have 4 years of data and it takes about 1:30 minutes to import all the data)

After the import, you can query against the sqlite3 database in `store/athletes/<id>.db` using some [3rd party tool](https://www.sqliteflow.com/), org-mode, or the command line application `sqlite3`.

## Plan App

Plan takes in a [specifically formatted `.CSV` file](https://github.com/robrohan/hungrylegs/tree/master/cmd/plan) and makes `.ICS` file that you can use on your desktop or google calendar.
