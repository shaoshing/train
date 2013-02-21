require "sass"
require "socket"
require File.expand_path('../lib/interpreter', __FILE__)

SOCKET_NAME = "/tmp/train.sass.socket"

class SassInterpreter < Interpreter
  def self.render content
    engine = Sass::Engine.new(content, :load_paths => ["assets/stylesheets"])
    engine.render
  end
end

SassInterpreter.run(SOCKET_NAME)
