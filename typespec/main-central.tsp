import "@typespec/http";
import "./models.tsp";

using Http;
using SharedModels;

@service(#{ title: "Report-Engine Central Service" })
@server("http://localhost:8080", "Report-Engine API (Central)")
@route("/api")
namespace CentralService;


@route("/ping")
@tag("ping") 
interface ping {
  @get op pong(): PingResponse;
}
