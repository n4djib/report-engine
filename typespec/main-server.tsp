import "@typespec/http";
import "./models.tsp";

using Http;
using SharedModels;

@service(#{ title: "Report-Engine Server Service" })
@server("http://localhost:8080", "Report-Engine API (Server)")
@route("/api")
namespace ServerService;


@route("/ping")
@tag("ping") 
interface ping {
  @get op pong(): PingResponse;
}
