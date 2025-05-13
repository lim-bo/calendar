#include "cfg.h"
QSettings settings("hihik", "hahak");
void load() {
    QSettings settings("hihik", "hahak");
    settings.setValue("host", "localhost");
    settings.setValue("port", "8080");
}
