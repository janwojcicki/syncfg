#include <iostream>
#include <cstring>
#include <unistd.h>
#include <sys/types.h>
#include <string>
#include <fstream>
#include <pwd.h>
#include <vector>
#include <filesystem>
#include <curlpp/cURLpp.hpp>
#include <curlpp/Easy.hpp>
#include <curlpp/Options.hpp>


using namespace curlpp::options;
struct passwd *pw = getpwuid(getuid());

const char *homedir = pw->pw_dir;
const std::string hdir (homedir);
std::vector<std::string> config;


void load_config(){
	std::ifstream inp (hdir+"/.kr");
	std::string line;
	while(std::getline(inp, line)){
		config.push_back(line);
	}
	inp.close();
}

void save_config(){
	std::ofstream inp (hdir+"/.kr");
	for (std::string s : config){
		inp << s << "\n";	
	}
	inp.close();
}

std::string expand_path(std::string s){
	if(s.substr(0,1) == "/")
		return s;
	else if(s.substr(0,1) == "~")
		return hdir+s.substr(1, s.size()-1);
	else 
		return std::string(std::filesystem::current_path())+"/"+s;
}

int main(int argc, char* argv[]){
	curlpp::Cleanup cleaner;
	curlpp::Easy request;

	std::list<std::string> headers;
	headers.push_back("Content-Type: multipart/form-data"); 
	request.setOpt(new curlpp::options::Url("http://127.0.0.1:5000/uploader")); 
	request.setOpt(new curlpp::options::HttpHeader(headers)); 
	curlpp::Forms formParts;
      formParts.push_back(new curlpp::FormParts::Content("file", "value1"));
	request.setOpt(new curlpp::options::HttpPost(formParts)); 
	request.perform();



	std::vector<std::string> args;
	for (int i = 0; i < argc; i++){
		args.push_back(std::string(argv[i]));
	}

	if(argc > 1){
		if (args[1] == "add" && argc == 5){
			load_config();
			bool found = false;
			for (uint j = 0; j < config.size(); j++){
				std::string s = config[j];
				if (s.substr(0,1) == "#"){
					if (s.substr(1, s.size()-1) == args[2])	{
						found = true;
						config.insert(config.begin()+j+1, args[3]+":"+expand_path(args[4]));
						break;
					}
				}
			}
			if(!found){
				config.push_back("#"+args[2]);
				config.push_back(args[3]+":"+expand_path(args[4]));
			}
			save_config();
			return 0;
		}	
	}	
	return 0;
}
