#!/usr/bin/env ruby

exists = `gem list | grep awesome_print`
puts `gem install awesome_print` if exists.empty?

require 'time'
require 'awesome_print'

ap 'Starting auto_backup..'

MIN_H, MAX_H = 6, 22
BACK_DATE = 2

now = Time.now

class Time
	def iso
		self.strftime('%F %T')
	end
end

def del_backup what
	num = `ls -w 1 /home/CHANGEME/backup/full_backup--#{what}*.sql.xz | wc -l`
	return if num.strip.to_i == 0
	ap "Remove old backup: #{what}*"
	puts `rm /home/CHANGEME/backup/full_backup--#{what}*.sql.xz`
end

with_log = ''

loop do
	loop do
		with_log = ''
		delta = 1
		delta = 2 if now.saturday?
		delta = 3 if now.sunday?
		new_hour = now.hour+delta
		old_day = now.day
		now = now + delta * 3600 # increment x hour
		if new_hour > MAX_H and old_day == now.day # still on the same day but greater than 22:00
			now = Time.new(now.year, now.month, now.day, 23, 45)
			with_log = '_with_log'
			break
		end
		break if (now.hour >= MIN_H and now.hour <= MAX_H)
	end
	dur = (now - Time.now).to_i + 1
	ap "#{Time.now.iso}, sleep #{dur} seconds until #{now.iso}"
	puts 'INCORRECT TIME !!!' if dur < 1
	dur = 1 if dur < 1
	sleep dur
	puts `./full_backup#{with_log}.sh`
	del_backup (Time.now-60*60*24*BACK_DATE).strftime('%F')
end
