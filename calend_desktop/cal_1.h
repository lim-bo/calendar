#ifndef CAL_1_H
#define CAL_1_H
#include <QHBoxLayout>
#include <QMessageBox>
#include <QWidget>
#include "event_entry.h"
#include "client.h"
#include "eventdata.h"

namespace Ui {
class cal_1;
}

class cal_1 : public QWidget
{
    Q_OBJECT

public:
    explicit cal_1(QString uid, QWidget *parent = nullptr);
    ~cal_1();
private slots:
    void loadEventsForToday();
    void eventDeleted(event_entry* ev);
    void showEvent(QShowEvent *event);
    void eventEdited(event_entry* ev);

private:
    Ui::cal_1 *ui;
    Client cli;
    QString uid;
};

#endif // CAL_1_H
