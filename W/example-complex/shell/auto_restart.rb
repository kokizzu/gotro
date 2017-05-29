#!/usr/bin/env ruby

exists = `gem list | grep awesome_print`
puts `gem install awesome_print` if exists.empty?

require 'time'
require 'awesome_print'

now = Time.now

class Time
	def iso
		self.strftime('%F %T')
	end
end

loop do
	now = Time.new(now.year, now.month, now.day, 24)
	dur = (now - Time.now).to_i + 1
	ap "#{Time.now.iso}, sleep #{dur} seconds until #{now.iso}"
	puts 'INCORRECT TIME !!!' if dur < 1
	dur = 1 if dur < 1
	sleep dur
	puts `systemctl restart CHANGEME`
end
