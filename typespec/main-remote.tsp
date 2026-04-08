import "@typespec/http";
import "./models.tsp";

using Http;
using SharedModels;

@service(#{ title: "Report-Engine Remote Service" })
// TODO make this configurable
@server("http://localhost:8081", "Report-Engine API (Remote)")
@route("/api")
namespace RemoteService;


@route("/ping")
@tag("ping")
interface ping {
  @get op pong(): PingResponse;
}
