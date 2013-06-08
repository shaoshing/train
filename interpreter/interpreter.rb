require "socket"

class Interpreter
  SOCKET_NAME = ARGV[0]
  MASTER_PID = ARGV[1].to_i

  def self.serve
    prepare_for_automatic_termination

    server = listen
    loop do
      client = server.accept
      format, option, content, file_path = read_all(client)

      begin
        result = self.send("render_#{format}", content, option, file_path)
        client.write "success<<#{result}"
      rescue Exception => e
        puts e
        client.write "error<<#{e}"
      end
      client.close
    end
  end

  # Shut down the interpreter when master process does not exist.
  def self.prepare_for_automatic_termination
    interpreter_pid = Process.pid
    fork do
      $0 = "ruby [train] master process monitor (PID #{MASTER_PID})" # Set Process Name
      loop do
        sleep 2
        begin
          Process.getpgid(MASTER_PID)
        rescue #=> when master process cannot be found
          Interpreter.clean_up
          Process.kill 1, interpreter_pid
          exit
        end
      end
    end
  end

  def self.clean_up
    `rm -f #{SOCKET_NAME}`
  end

  private

  def self.listen
    begin
      clean_up

      server = UNIXServer.new(SOCKET_NAME)
      puts "<<ready"
      `touch #{SOCKET_NAME}`
      server
    rescue => e
      puts e
      exit 1
    end
  end

  def self.read_all client
    data = ""
    recv_length = 2000
    while tmp = client.recv(recv_length)
      data += tmp
      break if tmp.length < recv_length
    end
   data.split("<<")
  end

  def self.render_sass content, option, file_path
    _render_sass(content, :sass, option, file_path)
  end

  def self.render_scss content, option, file_path
    _render_sass(content, :scss, option, file_path)
  end

  def self.render_sass_source_map content, option, file_path
    _render_sass(content, :scss, "source_map", file_path)
  end

  def self._render_sass content, syntax, option, file_path
    require "sass"

    options = {
      :load_paths => ["assets/stylesheets"],
      :syntax => syntax
    }

    options[:debug_info] = true if option == "debug_info"
    options[:line_numbers] = true if option == "line_numbers"
    options[:filename] = file_path
    engine = Sass::Engine.new(content, options)

    css_uri = file_path.sub(/\.(scss|sass)/, ".css")
    @source_map ||= {}
    return @source_map[css_uri] if option == "source_map" && @source_map[css_uri]

    if engine.respond_to?(:render_with_sourcemap)
      source_map_uri = "/" + file_path + ".map"
      results = engine.render_with_sourcemap(source_map_uri)
      css = results[0]
      source_map = results[1].to_json(:css_uri => css_uri, :sourcemap_path => source_map_uri)

      @source_map[css_uri] = source_map
      option == "source_map" ? @source_map[css_uri] : css
    else
      if option == "source_map"
        raise "Please install sass pre-released SASS version (gem in sass --pre) to get support of sourcemap"
      end

      engine.render
    end
  end

  def self.render_coffee content, option, file_path
    require "coffee-script"
    CoffeeScript.compile content
  end
end

trap("INT"){} #=> make silent the "Interrupted" error
Interpreter.serve
